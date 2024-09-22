package sportspace

type Config struct {
	OTPLength uint `env:"OTP_LENGTH" envDefault:"6"`
}
