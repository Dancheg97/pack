package tmpl

const Gnupgerr = `GPG key is not found in user directory ~/.gnupg
It is required for package signing, run:

1) Install gnupg:
pack i gnupg

2) Generate a key:
gpg --gen-key

3) Get KEY-ID, paste it to next command:
gpg -k

4) Send it to key server:
gpg --send-keys KEY-ID`
