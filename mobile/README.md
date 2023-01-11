To build an Android Archive you will need to have a suitable Android NDK installed and `ANDROID_NDK_HOME` env variable
set, e.g.

```bash
export ANDROID_NDK_HOME=/Users/rowan/Library/Android/sdk/ndk/22.1.7171670
```

To create the AAR run:

```bash
go get golang.org/x/mobile/cmd/gomobile
gomobile init
gomobile bind -target android -javapkg=com.nyaruka.goflow -o mobile/goflow.aar github.com/nyaruka/goflow/mobile
```