package handler_test

var RsaPrivateKey string = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

var RsaPublicKey string = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAzDwd6/ye1D07XGzED5hE
2Hcwhj+kZypEweGcHSZxAcgz1+sZtT560bZk1IpOpN6u3zIa3lXH7VSJcppw2CJS
c/wQRzz2/PXCG4hOIo1A0gpzYXIVo4qjCOM7v9ObjzbVo/9e0+VISwrRjcZpAeah
NAlGOr8KtOuV/o4eIDlxxDEuGwTtzJ4L20CniM2BFqVwOE/aSuJldQAAghIb7od5
qeh3U3gzFRDzbyEMQMMxHEkdpkB4DHlmgpeth1ptwKKAbXdFR/Pd3axSiCrGDYd1
3CI5Ad607vBkajEuXyi5PVWPRHHlRLKcN/qIyYCaukqtavEWrOCjAAefmhNkOT10
S584PH7iofPYAxCE+f8hPQaiUvVJ9IVZvqU14f8GU7W2dSpkFOfJ1Z2sEHdX7deb
A8XWkJmmXPDKGBz03+a22m1aU9O1L3xHz4SrO9XGAHZMR+Vbq5mVSwiXKncyYJ7a
1vEB3b5ulkJXDPZljw4U6l3KINM80cg9laC2lv1ejALA4O+/cl+5iXHXU8vSluvY
Uk3iWyXYAph21plmuNPhDijty3kaXy9pN4J0GLFWGhX+VbU13mI/d761V1qtgaQh
Te0TTMcUWug38JK7+mZok8pSk9Mg367KoOnuYuLzcG8IMwbNAh6FpwVR2stl58aK
L20KPET5bNpzqtYUqdM/K/ECAwEAAQ==
-----END PUBLIC KEY-----
`

var dummyJWT string = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImMxMThhMWE5LTI4ZjEtNDEzNy05MDkzLTg3NDg3ZDI0ZTVkOSIsImZ1bGxfbmFtZSI6IlRlc3QgVXNlciIsInBob25lX251bWJlciI6Iis2MjgxMjM0NTY3ODkiLCJleHAiOjE5Njg5MDg0NDR9.B0EI5PaAFzoQ0tlvobzYePW_NJ6cFxF_T-Kr5y5KNiUxOUg3XRXzk7sWWBSyzYFuK-tBOMDLxpNVAfOIL0xqiVVrDm21miEVFeVU1-AjzGLJZ_rkraIkSx-J0qFLi5cGHQo8D4TGiuWhANDNfbooAo66DzAbVQ2oxggefaQrxMUWpejNLKOkyJrSEKudFazMvDvXymHtssLNrLmeJeieozHtbHmBplMbx78GM0QXBiOkYmYCGetbmaHTBE9-2mn55Kkq_JLMwUx3s4EzUcEjk1s9AX-MNSQAadpLmnOGQKqFmCqX6s7aSk5l3dy9jPof_F1_0iwdevtDM4u690FMdXWjnrpOLtIiRc5abcnSxW3Dl2b3xv7NyEqvz7d2hRRfEIQYJfcZsxd5hXNCJ2QiUj5cyX0yTYWxjHTaTCgpb7ZKpkm1zxxRIhvSnj1y4U2wIfSvAc2OHcWLuGRsgCXlgW4ejoHgQD2X7DsSlVgF8AGSUW8JD159x_LOFJ01FsK99ahzFRsocKMJW_csCDuWRZivILwmaLbx-47Tq_fEcI1mln4IrVwskfXy9aBHclYgOuQ0hTsFFx4ZIBoFeFz2ib2839pAjDbIqa-kuFS3h199mLB7pnr_PEgvm4tHaKd05vc5OpkjgY8AUbFA3uLtZ8GPhAhliVpEBaN7-nKxl3g"
