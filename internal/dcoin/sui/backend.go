package sui

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/ios1024/dpass/internal/dcoin"
	"github.com/ios1024/dpass/pkg/bech32"
	"github.com/ios1024/dpass/pkg/bip3x"
	"github.com/ios1024/dpass/pkg/hashx"
	"github.com/spf13/cobra"
)

const (
	purposeDefault = 44
	coinDefault    = 784
	accountDefault = 0
	changeDefault  = 0
	indexDefault   = 0
	secretDefault  = false
)

var (
	errInvalidPurpose = errors.New("invalid purpose")
	errInvalidCoin    = errors.New("invalid coin")
	errInvalidAccount = errors.New("invalid account")
	errInvalidChange  = errors.New("invalid change")
	errInvalidIndex   = errors.New("invalid index")
)

type backend struct {
	purpose uint32
	coin    uint32
	account uint32
	change  uint32
	index   uint32
	secret  bool
}

func backendDefault() *backend {
	return &backend{
		purpose: purposeDefault,
		coin:    coinDefault,
		account: accountDefault,
		change:  changeDefault,
		index:   indexDefault,
		secret:  secretDefault,
	}
}

func NewCmd() *cobra.Command {
	backend := backendDefault()
	cmd := &cobra.Command{Use: "sui", Args: cobra.NoArgs, RunE: backend.runE}
	cmd.Flags().Uint32Var(&backend.index, "index", indexDefault, fmt.Sprintf(
		"index is the number of address (default %v)", indexDefault))
	cmd.Flags().BoolVar(&backend.secret, "secret", secretDefault, fmt.Sprintf(
		"show secret instead of address (default %t)", secretDefault))
	return cmd
}

func (b *backend) checkArguments() error {
	if b.purpose >= bip3x.FirstHardenedChild {
		return errInvalidPurpose
	}
	if b.coin >= bip3x.FirstHardenedChild {
		return errInvalidCoin
	}
	if b.account >= bip3x.FirstHardenedChild {
		return errInvalidAccount
	}
	if b.change >= bip3x.FirstHardenedChild {
		return errInvalidChange
	}
	if b.index >= bip3x.FirstHardenedChild {
		return errInvalidIndex
	}
	return nil
}

func (b *backend) runE(_ *cobra.Command, _ []string) error {
	mnemonic, err := dcoin.ReadMnemonic()
	if err != nil {
		return fmt.Errorf("failed to read mnemonic: %w", err)
	}
	result, err := b.getResult(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to get result: %w", err)
	}
	if _, err := os.Stdout.WriteString(result); err != nil {
		return fmt.Errorf("failed to write result: %w", err)
	}
	return nil
}

func (b *backend) getResult(mnemonic string) (string, error) {
	if err := b.checkArguments(); err != nil {
		return "", fmt.Errorf("failed to check arguments: %w", err)
	}
	seed, err := bip3x.MnemonicToSeed(mnemonic, "")
	if err != nil {
		return "", fmt.Errorf("failed to convert mnemonic to seed: %w", err)
	}
	sk, err := bip3x.Ed25519DeriveSk(seed, []uint32{
		b.purpose + bip3x.FirstHardenedChild,
		b.coin + bip3x.FirstHardenedChild,
		b.account + bip3x.FirstHardenedChild,
		b.change + bip3x.FirstHardenedChild,
		b.index + bip3x.FirstHardenedChild,
	})
	if err != nil {
		return "", fmt.Errorf("failed to derive sk: %w", err)
	}
	privateKey := ed25519.NewKeyFromSeed(sk)
	if b.secret {
		return skToWIF(privateKey[:ed25519.SeedSize]), nil
	}
	return pkToAddress(privateKey[ed25519.SeedSize:]), nil
}

func skToWIF(sk []byte) string {
	return bech32.Encode("suiprivkey", nil, slices.Concat([]byte{0}, sk))
}

func pkToAddress(pk []byte) string {
	return "0x" + hex.EncodeToString(hashx.Blake2b256Sum(slices.Concat([]byte{0}, pk)))
}
