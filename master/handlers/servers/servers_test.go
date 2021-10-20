package servers

import "testing"

func TestGenSubscriptionData(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				uid: "hello123",
			},
			want:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterServer("hello.lazarus", "hello2.lazarus", "thisislazarus")
			got, err := GenSubscriptionData(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenSubscriptionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenSubscriptionData() = %v, want %v", got, tt.want)
			}
		})
	}
}
