// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: Copyright 2024 Stacklok, Inc

package mindpack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestManifestWrite(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		manifest *Manifest
		mustErr  bool
	}{
		{
			"normal",
			&Manifest{
				Metadata: &Metadata{
					Name:      "test",
					Namespace: "testspace",
					Version:   "v1.2.0",
					Date: &timestamppb.Timestamp{
						Seconds: 1709866805,
					},
				},
				Files: &Files{
					Profiles: []*File{
						{
							Name: "profile.yaml",
							Hashes: map[string]string{
								SHA256: "8b438ca800dfa20c6ca66ed83f05ef874cc1e1859d1a0a193b4c0727e5629977",
							},
						},
					},
					RuleTypes: []*File{
						{
							Name: "rule_type.yaml",
							Hashes: map[string]string{
								SHA256: "0aecaf4d7ce19dc39679952c6951005e1396a5e615289ff3deb351873957d055",
							},
						},
					},
				},
			},
			false,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			err := tc.manifest.Write(b)
			if tc.mustErr {
				require.Error(t, err)
				return
			}
			man := &Manifest{}
			j := b.Bytes()
			fmt.Printf("%s", string(j))
			require.NoError(t, json.Unmarshal(b.Bytes(), man))

			if diff := cmp.Diff(tc.manifest, man, protocmp.Transform()); diff != "" {
				t.Fatalf("assertion failed: values are not equal\n--- expected\n+++ actual\n%v", diff)
			}
		})
	}

}
