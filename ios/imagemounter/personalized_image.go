package imagemounter

import (
	"fmt"
	"howett.net/plist"
	"os"
	"strconv"
	"strings"
)

type buildManifest struct {
	BuildIdentities []buildIdentity
}

func loadBuildManifest(p string) (buildManifest, error) {
	f, err := os.Open(p)
	if err != nil {
		return buildManifest{}, err
	}
	defer f.Close()
	dec := plist.NewDecoder(f)
	var m buildManifest
	err = dec.Decode(&m)
	if err != nil {
		return buildManifest{}, err
	}
	return m, nil
}

func (m buildManifest) findIdentity(identifiers personalizationIdentifiers) (buildIdentity, error) {
	for _, i := range m.BuildIdentities {
		if i.ApBoardID() == identifiers.BoardId && i.ApChipID() == identifiers.ChipID {
			return i, nil
		}
	}
	return buildIdentity{}, fmt.Errorf("not found")
}

type buildIdentity struct {
	BoardID  string `plist:"ApBoardID"`
	ChipID   string `plist:"ApChipID"`
	Manifest struct {
		LoadableTrustCache struct {
			Digest []byte
		}
		PersonalizedDmg struct {
			Digest []byte
		} `plist:"PersonalizedDMG"`
	}
}

func (b buildIdentity) ApBoardID() int {
	return hexToInt(b.BoardID)
}

func (b buildIdentity) ApChipID() int {
	return hexToInt(b.ChipID)
}

type personalizationIdentifiers struct {
	BoardId        int
	ChipID         int
	SecurityDomain int
}

func hexToInt(s string) int {
	i, err := strconv.ParseInt(strings.ReplaceAll(strings.ToLower(s), "0x", ""), 16, 32)
	if err != nil {
		return 0
	}
	return int(i)
}
