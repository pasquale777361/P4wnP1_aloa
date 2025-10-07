package cli_client

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "github.com/mame82/P4wnP1_aloa/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/status"
)

var (
	btAlias        string
	btDiscoverable bool
	btPairable     bool
	btSSP          bool
	btScan         bool
)

var bluetoothCmd = &cobra.Command{
	Use:   "bluetooth",
	Short: "Configure Bluetooth settings",
	Long:  `Configure Bluetooth settings, such as alias, discoverability, and pairability.`,
}

var bluetoothGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current Bluetooth settings",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn, err := Client(StrRemoteHost, StrRemotePort)
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.Empty{}
		settings, err := client.BluetoothGetSettings(ctx, req)
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}
		fmt.Printf("Bluetooth Settings:\n")
		fmt.Printf("  Alias: %s\n", settings.Ci.Name)
		fmt.Printf("  Discoverable: %t\n", settings.Ci.CurrentSettings.Discoverable)
		fmt.Printf("  Pairable: %t\n", settings.Ci.CurrentSettings.Bondable)
		fmt.Printf("  SSP: %t\n", settings.Ci.CurrentSettings.SecureSimplePairing)
		fmt.Printf("  Powered: %t\n", settings.Ci.CurrentSettings.Powered)
	},
}

var bluetoothSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set Bluetooth controller settings",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn, err := Client(StrRemoteHost, StrRemotePort)
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		currentSettings, err := client.BluetoothGetSettings(ctx, &pb.Empty{})
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}

		if cmd.Flags().Changed("alias") {
			currentSettings.Ci.Name = btAlias
		}
		if cmd.Flags().Changed("discoverable") {
			currentSettings.Ci.CurrentSettings.Discoverable = btDiscoverable
		}
		if cmd.Flags().Changed("pairable") {
			currentSettings.Ci.CurrentSettings.Bondable = btPairable
		}
		if cmd.Flags().Changed("ssp") {
			currentSettings.Ci.CurrentSettings.SecureSimplePairing = btSSP
		}

		_, err = client.BluetoothSetSettings(ctx, currentSettings)
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}
		fmt.Println("Bluetooth settings updated successfully.")
	},
}

var bluetoothOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Power on the Bluetooth controller",
	Run: func(cmd *cobra.Command, args []string) {
		setBluetoothPower(true)
	},
}

var bluetoothOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Power off the Bluetooth controller",
	Run: func(cmd *cobra.Command, args []string) {
		setBluetoothPower(false)
	},
}

func setBluetoothPower(powerState bool) {
	client, conn, err := Client(StrRemoteHost, StrRemotePort)
	if err != nil {
		fmt.Println(status.Convert(err).Message())
		os.Exit(-1)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	settings, err := client.BluetoothGetSettings(ctx, &pb.Empty{})
	if err != nil {
		fmt.Println(status.Convert(err).Message())
		os.Exit(-1)
	}

	settings.Ci.CurrentSettings.Powered = powerState

	_, err = client.BluetoothSetSettings(ctx, settings)
	if err != nil {
		fmt.Println(status.Convert(err).Message())
		os.Exit(-1)
	}

	if powerState {
		fmt.Println("Bluetooth controller powered on.")
	} else {
		fmt.Println("Bluetooth controller powered off.")
	}
}

var bluetoothScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for nearby Bluetooth devices",
	Run: func(cmd *cobra.Command, args []string) {
		client, conn, err := Client(StrRemoteHost, StrRemotePort)
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}
		defer conn.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		stream, err := client.BluetoothScan(ctx, &pb.Empty{})
		if err != nil {
			fmt.Println(status.Convert(err).Message())
			os.Exit(-1)
		}

		fmt.Println("Scanning for Bluetooth devices...")
		for {
			device, err := stream.Recv()
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				fmt.Println(status.Convert(err).Message())
				break
			}
			fmt.Printf("  - Address: %s, Name: %s, RSSI: %d\n", device.Address, device.Name, device.Rssi)
		}
		fmt.Println("Scan complete.")
	},
}

func init() {
	rootCmd.AddCommand(bluetoothCmd)
	bluetoothCmd.AddCommand(bluetoothGetCmd)
	bluetoothCmd.AddCommand(bluetoothSetCmd)
	bluetoothCmd.AddCommand(bluetoothOnCmd)
	bluetoothCmd.AddCommand(bluetoothOffCmd)
	bluetoothCmd.AddCommand(bluetoothScanCmd)

	bluetoothSetCmd.Flags().StringVar(&btAlias, "alias", "", "Set the Bluetooth device alias")
	bluetoothSetCmd.Flags().BoolVar(&btDiscoverable, "discoverable", false, "Set discoverable mode")
	bluetoothSetCmd.Flags().BoolVar(&btPairable, "pairable", false, "Set pairable mode")
	bluetoothSetCmd.Flags().BoolVar(&btSSP, "ssp", false, "Enable/disable Secure Simple Pairing")
}