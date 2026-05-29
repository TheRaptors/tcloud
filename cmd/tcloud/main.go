package main

import (
	"fmt"
	"os"

	"my-project-golang/internal/config"
	"my-project-golang/internal/cvm"
	"my-project-golang/internal/dnspod"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置文件失败: %s\n", err)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "list":
		// 获取记录列表，提取 RecordId
		recordId, err := dnspod.DescribeRecordList(cfg)
		if err != nil {
			fmt.Printf("获取记录列表失败: %s\n", err)
			os.Exit(1)
		}
		// 可选：直接用获取到的 RecordId 查询详情
		if len(os.Args) > 2 && os.Args[2] == "--detail" {
			fmt.Println()
			if err := dnspod.DescribeRecord(cfg, recordId); err != nil {
				fmt.Printf("获取记录详情失败: %s\n", err)
				os.Exit(1)
			}
		}

	case "describe":
		// 先获取记录列表拿到 RecordId，再查询详情
		recordId, err := dnspod.DescribeRecordList(cfg)
		if err != nil {
			fmt.Printf("获取记录列表失败: %s\n", err)
			os.Exit(1)
		}
		fmt.Println()
		if err := dnspod.DescribeRecord(cfg, recordId); err != nil {
			fmt.Printf("获取记录详情失败: %s\n", err)
			os.Exit(1)
		}

	case "modify":
		// 修改记录：先获取 RecordId，再修改
		if len(os.Args) < 3 {
			fmt.Println("用法: tcloud modify <新IP地址>")
			fmt.Println("示例: tcloud modify 200.200.200.200")
			os.Exit(1)
		}
		value := os.Args[2]
		recordId, err := dnspod.DescribeRecordList(cfg)
		if err != nil {
			fmt.Printf("获取RecordId失败: %s\n", err)
			os.Exit(1)
		}
		fmt.Println()
		if err := dnspod.ModifyRecord(cfg, recordId, value); err != nil {
			fmt.Printf("修改记录失败: %s\n", err)
			os.Exit(1)
		}

	case "run-instances":
		// 创建CVM实例
		instanceId, err := cvm.RunInstances(cfg)
		if err != nil {
			fmt.Printf("创建实例失败: %s\n", err)
			os.Exit(1)
		}
		_ = instanceId

	case "deploy":
		// 一键部署：购买竞价实例 → 获取公网IP → 修改DNS记录
		fmt.Println("========== 第1步：购买竞价实例 ==========")
		instanceId, err := cvm.RunInstances(cfg)
		if err != nil {
			fmt.Printf("创建实例失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第2步：获取公网IP ==========")
		publicIP, err := cvm.DescribeInstances(cfg, instanceId)
		if err != nil {
			fmt.Printf("获取公网IP失败: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("公网IP: %s\n", publicIP)

		fmt.Println("\n========== 第3步：获取 RecordId ==========")
		recordId, err := dnspod.DescribeRecordList(cfg)
		if err != nil {
			fmt.Printf("获取RecordId失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第4步：修改前的域名解析信息 ==========")
		if err := dnspod.DescribeRecord(cfg, recordId); err != nil {
			fmt.Printf("查询修改前记录失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第5步：修改DNS A记录 ==========")
		fmt.Printf("将 %s.%s 指向 %s\n", cfg.Subdomain, cfg.Domain, publicIP)
		if err := dnspod.ModifyRecord(cfg, recordId, publicIP); err != nil {
			fmt.Printf("修改记录失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第6步：修改后的域名解析信息 ==========")
		if err := dnspod.DescribeRecord(cfg, recordId); err != nil {
			fmt.Printf("查询修改后记录失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 部署完成 ==========")
		fmt.Printf("实例ID: %s\n", instanceId)
		fmt.Printf("公网IP: %s\n", publicIP)
		fmt.Printf("域名: %s.%s → %s\n", cfg.Subdomain, cfg.Domain, publicIP)

	case "destroy":
		// 根据内网IP查找并销毁CVM实例
		fmt.Println("========== 根据内网IP查找实例 ==========")
		instanceId, err := cvm.FindInstanceByPrivateIP(cfg)
		if err != nil {
			fmt.Printf("查找实例失败: %s\n", err)
			os.Exit(1)
		}
		fmt.Println()
		if err := cvm.TerminateInstances(cfg, instanceId); err != nil {
			fmt.Printf("销毁实例失败: %s\n", err)
			os.Exit(1)
		}

	case "undeploy":
		// 一键回收：查找实例→销毁实例→还原DNS记录
		fmt.Println("========== 第1步：根据内网IP查找实例 ==========")
		instanceId, err := cvm.FindInstanceByPrivateIP(cfg)
		if err != nil {
			fmt.Printf("查找实例失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第2步：获取 RecordId ==========")
		recordId, err := dnspod.DescribeRecordList(cfg)
		if err != nil {
			fmt.Printf("获取RecordId失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第3步：销毁前的域名解析信息 ==========")
		if err := dnspod.DescribeRecord(cfg, recordId); err != nil {
			fmt.Printf("查询销毁前记录失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第4步：销毁CVM实例 ==========")
		if err := cvm.TerminateInstances(cfg, instanceId); err != nil {
			fmt.Printf("销毁实例失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第5步：修改DNS A记录（还原为 0.0.0.0） ==========")
		fmt.Printf("将 %s.%s 指向 0.0.0.0\n", cfg.Subdomain, cfg.Domain)
		if err := dnspod.ModifyRecord(cfg, recordId, "0.0.0.0"); err != nil {
			fmt.Printf("修改记录失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 第6步：修改后的域名解析信息 ==========")
		if err := dnspod.DescribeRecord(cfg, recordId); err != nil {
			fmt.Printf("查询修改后记录失败: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("\n========== 回收完成 ==========")
		fmt.Printf("已销毁实例: %s\n", instanceId)
		fmt.Printf("域名: %s.%s → 0.0.0.0\n", cfg.Subdomain, cfg.Domain)

	default:
		fmt.Printf("未知命令: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("腾讯云 DNSPod & CVM 管理工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  tcloud <命令> [参数]")
	fmt.Println()
	fmt.Println("命令:")
	fmt.Println("  list              获取DNS解析记录列表（提取RecordId）")
	fmt.Println("  list --detail     获取记录列表后，自动查询第一条记录详情")
	fmt.Println("  describe          自动获取RecordId并查询记录详情")
	fmt.Println("  modify <IP>       自动获取RecordId并修改记录的IP值")
	fmt.Println("  run-instances     创建CVM竞价实例")
	fmt.Println("  deploy            一键部署：购买实例→获取公网IP→修改DNS")
	fmt.Println("  destroy           根据内网IP自动查找并销毁CVM实例")
	fmt.Println("  undeploy          一键回收：自动查找实例→销毁→还原DNS")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run ./cmd/tcloud list")
	fmt.Println("  go run ./cmd/tcloud deploy")
	fmt.Println("  go run ./cmd/tcloud undeploy")
}
