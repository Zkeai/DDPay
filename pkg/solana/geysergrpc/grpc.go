package geysergrpc

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"github.com/Zkeai/go_template/pkg/solana/geysergrpc/proto"
	"github.com/mr-tron/base58"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net/url"
	"strings"
	"time"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second,
	Timeout:             60 * time.Second,
	PermitWithoutStream: true,
}

type GeyserClient struct {
	conn   *grpc.ClientConn
	client proto.GeyserClient
}

type GrpcConfig struct {
	Endpoint string `yaml:"Endpoint"`
}

var insecureConnection = false

func grpcConnect(address string, plaintext bool) *grpc.ClientConn {
	var opts []grpc.DialOption

	if plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	log.Println("启动 gRPC 客户端，连接到", address)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	return conn
}

func grpcSubscribe(conn *grpc.ClientConn, cfg SubscribeConfig) {
	client := proto.NewGeyserClient(conn)

	subscription := proto.SubscribeRequest{
		Transactions: map[string]*proto.SubscribeRequestFilterTransactions{
			"transactions_sub": {
				Failed:         cfg.FailedTransactions,
				Vote:           cfg.VoteTransactions,
				AccountInclude: cfg.TransactionsAccountsInclude,
				AccountExclude: cfg.TransactionsAccountsExclude,
			},
		},
	}

	// 仅订阅指定签名（如果有）
	if cfg.Signature != nil && *cfg.Signature != "" {
		tr := true
		subscription.Transactions["signature_sub"] = &proto.SubscribeRequestFilterTransactions{
			Failed:    &tr,
			Vote:      &tr,
			Signature: cfg.Signature,
		}
	}

	subscriptionJson, err := json.Marshal(&subscription)
	if err != nil {
		log.Printf("序列化订阅请求失败: %v", err)
	}
	log.Printf("订阅请求: %s", string(subscriptionJson))

	ctx := context.Background()
	if cfg.Token != "" {
		md := metadata.New(map[string]string{"x-token": cfg.Token})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("订阅创建失败: %v", err)
	}
	err = stream.Send(&subscription)
	if err != nil {
		log.Fatalf("发送订阅失败: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("订阅结束: %v", err)
			return
		}
		if err != nil {
			log.Fatalf("接收订阅失败: %v", err)
		}

		if tx := resp.GetTransaction(); tx != nil {
			accountKeys := tx.GetTransaction().Transaction.GetMessage().GetAccountKeys()
			accountKeyStrs := make([]string, len(accountKeys))
			for i, key := range accountKeys {
				accountKeyStrs[i] = base58.Encode(key)
			}

			// 识别命中的订阅地址
			matchedAccounts := []string{}
			includeMap := make(map[string]bool)
			for _, addr := range cfg.TransactionsAccountsInclude {
				includeMap[addr] = true
			}
			for _, addr := range accountKeyStrs {
				if includeMap[addr] {
					matchedAccounts = append(matchedAccounts, addr)
				}
			}

			for _, logMessage := range tx.GetTransaction().Meta.GetLogMessages() {
				if strings.Contains(logMessage, "Program log: Instruction: Buy") || strings.Contains(logMessage, "Program log: Instruction: Sell") {
					for _, msg := range tx.GetTransaction().Meta.GetLogMessages() {
						if strings.Contains(msg, "Program data: ") {
							data := strings.Split(msg, "Program data: ")[1]
							decoded, err := base64.StdEncoding.DecodeString(data)
							if err != nil {
								log.Printf("解码失败: %v", err)
								continue
							}

							// 验证数据长度
							if len(decoded) < 8+32+8+8+1+32+8 {
								continue
							}

							offset := 8 // 跳过魔数和版本

							// 解析Mint地址
							var mintBytes [32]byte
							copy(mintBytes[:], decoded[offset:offset+32])
							mintAddress := base58.Encode(mintBytes[:])
							offset += 32

							// 跳过sol和token数量
							offset += 16 // 8+8

							// 跳过isBuy标志
							offset += 1

							// 解析用户地址
							var userBytes [32]byte
							copy(userBytes[:], decoded[offset:offset+32])
							userAddress := base58.Encode(userBytes[:])
							offset += 32

							// 获取本地时间戳（毫秒）
							milliseconds := time.Now().UnixMilli()
							tradeTime := time.Now().Format("2006-01-02 15:04:05.000")

							// 输出交易信息
							log.Printf("\n===================== 买入交易 =====================\n"+
								"用户地址: %s\n"+
								"Mint地址: %s\n"+
								"交易时间: %s\n"+
								"时间戳(毫秒): %d\n"+
								"匹配地址: %s\n"+
								"================================================\n",
								userAddress,
								mintAddress,
								tradeTime,
								milliseconds,
								matchedAccounts,
							)
						}
					}
				}
			}
		}
	}
}

func NewGeyserClient(cfg SubscribeConfig) {
	log.Printf("📡 endpoint: %s", cfg.Endpoint)

	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		log.Fatalf("提供的 GRPC 地址无效: %v", err)
	}
	if u.Scheme == "http" {
		insecureConnection = true
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		if insecureConnection {
			port = "80"
		} else {
			port = "443"
		}
	}
	address := host + ":" + port

	conn := grpcConnect(address, insecureConnection)
	defer conn.Close()

	grpcSubscribe(conn, cfg)
}

func decodeUint64(b []byte) uint64 {
	if len(b) != 8 {
		return 0
	}
	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
}
