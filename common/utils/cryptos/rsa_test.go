package cryptos

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestRsaSignWithhash(t *testing.T) {
	priK := `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDf6lu4mo+yQziS0qEc1uoR5qAGwRvXwV0AnAIKKTG12/dHcqY37gtDCKmwjW7EdZAnkZTMDAzNlVsEtVP7XnzC65yY+mRTl3iG/rWMvUCNaUXBcZeW64XdVGRJoqCkvjd4jTFp1S99uQpG+Nb5gqfvlJGbjbIftF23Ei5epoTUv0NHQW6ulvQCTujF+cgIcUEjmyq7Dj+w8qWNcjDMHzvwjrC+G5sDbmdnOMXJZEBZjAyulBhADGHLLEK+YvJM51MKPK+a8SgR7/VY0+q0mefsihSHWvP99tyT61LRncX/o3m4VB1RNpYkGmWfK52zOraqCM0n3vAiHpBY6BZLO/krAgMBAAECggEBAJaRO+4NmNTxGMi58/a1mZ5B65e/IN7bOpOfVEvK2Y+Fg2k68gSoAFCqMZjz8ekPeMjyvxDahX10kki/OeLM7a4QyzOfI/mF4Fk+S4yA8jhk8rAalym35EMpbWqKfeVpt1lL8E9POGkdFYkV6VDMh+q1h5gHFyD0oxPomN+yr51yYTSpXAUFBOzMsD2bFNgxpof2Wu9tsp4hqKLLVfWRDla47rfzUCwh+Bp4PD21IbYNk3oAGrvmjuon4LPTDrJg0J+jS77A/eWpVIu/rQ8b7DdM8DhtmKopHens9uslSZuxRuIRpG2hR2rN0dcfJFYFsBqTVy/BWKpniJFvPjZ7RYECgYEA+AptSGCulZjzEC5tFc+rx6hb0CifgqjuA8dnpnDaHzo6pLsmVed8xERYaIEz/PJwvJ6M2Mk5cI7+9dSPA1h8OMNF8KXE9HxvIk1F3xN6UVOr84wYcX1kmoyMPfGsJOPcClcTwDuMbXL4hSLuxMVBOYv5slr6hKOqerDIypY0rOsCgYEA5xnAEE8ZK89N7t+2vVgzjtxvFOwn5ARZ2LgPFCpr4GvQUtqXNXCC9KQKB1rCCh5RPbhTS8eEIDtsSIGLd9BgO95R13FR4xEDryUFQqR+JX5C4NvbZfDAZ556SsnDfULGqLyVm94UUcjpzdIzT+GJhiGuaIN5X6mOPLkRg0o51MECgYAIeVq4cU0loT8Um3FwoFKvFIpmdyzT6u+Ow35ACnT5QiUEwbwSjUEO94LJtzhOeP3vA7+uHFnRBaGiRmvIYnqD+e/mw9MRwzqMwnUTpPe11ZT3Uh73qaAJQ6n658nIzNwUolrzY1Vt29KvwbzEjjSnQaf7Nu5+H5VQcb+6ZB1SlwKBgHyDJTYkR92Qzd577ktJ8E1yeu785ek2ZuobERS+Xm0F6bIaUAnc2tHQaA7aWV12RDNK0qYrkwaCva67DVe6j37yI4o+Ze4/RorhGVp0ofq1cncQPb9I3YF2o1EUMB2XEs3q/XiDSNNfuojIThkl3SDFmOB6pbRi+F3DIKpqHYqBAoGBAPVApWVDrf5Egl3+gYNAagzv2lzQ95L9a73xT+Cv+ILZyclO0ZgHJxWnSzHyYpMsz0XMa7jnh7osoLfPixYN8kB7DRonAofnyg7PePVA0IVR1U2D9tChQvAfO6iLUTB/ptlQB0iM+cZBQQSYZ6IWxGjXeWnALbAangINhUhSgnzR
-----END PRIVATE KEY-----`
	pk, err := ParsePriKeyPkcs8([]byte(priK))
	if err != nil {
		t.Error(err)
		return
	}
	data, err := RsaSignWithHash(pk, []byte(`{"Time":"1710899039","Rsa_ID":"Hxs9fa89"}`), 384)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(base64.StdEncoding.EncodeToString(data))
}
