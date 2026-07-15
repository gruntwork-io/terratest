package ssh_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh/v2"
	"github.com/stretchr/testify/require"
)

func TestSshConnectionOptions_ConnectionString(t *testing.T) {
	t.Parallel()

	type fields struct {
		Address string
		Port    int
	}

	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name: "plain ipv4",
			fields: fields{
				Address: "192.168.86.68",
				Port:    22,
			},
			want: "192.168.86.68:22",
		},
		{
			name: "plain ipv6",
			fields: fields{
				Address: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
				Port:    22,
			},
			want: "[2001:0db8:85a3:0000:0000:8a2e:0370:7334]:22",
		},
		{
			name: "host fqdn",
			fields: fields{
				Address: "host.for.test.com",
				Port:    443,
			},
			want: "host.for.test.com:443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			options := &ssh.SSHConnectionOptions{
				Address: tt.fields.Address,
				Port:    tt.fields.Port,
			}
			got := options.ConnectionString()
			require.Equal(t, tt.want, got)
		})
	}
}
