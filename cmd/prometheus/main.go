package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wisteriahuman/prometheus/internal/config"
	"github.com/wisteriahuman/prometheus/internal/server"
)

var version = "dev"

func main() {
	var port string

	rootCmd := &cobra.Command{
		Use:     "prm",
		Short:   "Markdownベースのノートアプリ Prometheus",
		Version: version,
	}

	devCmd := &cobra.Command{
		Use:   "dev [vault_path]",
		Short: "サーバーを起動",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vaultPath := ""
			if len(args) > 0 {
				vaultPath = args[0]
			}
			if port == "" {
				port = "3000"
			}

			cfg := config.New(vaultPath, "", port)
			fmt.Printf("Prometheus v%s\n", version)
			fmt.Printf("  vault: %s\n", cfg.VaultPath)
			fmt.Printf("  port:  %s\n", cfg.Port)
			fmt.Println()

			srv := server.NewServer(cfg)
			if err := srv.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
	devCmd.Flags().StringVarP(&port, "port", "p", "", "ポート番号 (default: 3000)")

	serveCmd := &cobra.Command{
		Use:   "serve [vault_path]",
		Short: "本番サーバーを起動",
		Args:  cobra.MaximumNArgs(1),
		Run:   devCmd.Run, // Same behavior for now
	}
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "ポート番号 (default: 3000)")

	initCmd := &cobra.Command{
		Use:   "init [vault_path]",
		Short: "新しいvaultを作成",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vaultPath := ""
			if len(args) > 0 {
				vaultPath = args[0]
			}
			cfg := config.New(vaultPath, "", "")
			fmt.Printf("Vault: %s\n", cfg.VaultPath)

			// Just create the vault directory — InitVault will be called when the server starts
			os.MkdirAll(cfg.VaultPath, 0o755)
			fmt.Println("✓ vault ディレクトリを作成しました")
			fmt.Println()
			fmt.Printf("次のコマンドで起動:\n  prm dev %s\n", cfg.VaultPath)
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "現在の設定を表示",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.New("", "", "")
			fmt.Println("Prometheus 設定情報")
			fmt.Println()
			fmt.Printf("  vault: %s\n", cfg.VaultPath)
			fmt.Printf("  db:    %s\n", cfg.DBPath)
			fmt.Printf("  port:  %s\n", cfg.Port)
			fmt.Println()
			fmt.Println("使い方:")
			fmt.Println("  prm dev ~/my-notes       サーバー起動")
			fmt.Println("  prm dev ~/work -p 3001   別ポートで起動")
			fmt.Println("  prm init ~/new-vault     vault作成")
		},
	}

	rootCmd.AddCommand(devCmd, serveCmd, initCmd, infoCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
