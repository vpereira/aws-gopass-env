# aws-gopass-env

A simple CLI tool to manage multiple AWS credential profiles using [gopass](https://www.gopass.pw/), and export them as environment variables for use with `aws-cli` and other AWS tools.

---

## Features

- Store AWS credentials securely in `gopass`
- Easily switch between multiple profiles
- Export credentials into your shell environment with `eval`
- Integration with `starship` prompt for automatic profile display
- Compatible with `zsh`, `bash`, and other POSIX shells

---

## Installation

### 1. Clone and build the project

```bash
git clone https://github.com/vpereira/aws-gopass-env.git
cd aws-gopass-env
go build -o aws-gopass-env
```

### 2. Setup a dedicated `gopass` store (optional but recommended)

This project expects credentials to live under a mounted store named `aws`. You can create one like this:

```bash
gopass mounts add aws ~/.password-store-aws
```

 **Note:** The store name `aws` is currently hardcoded. Support for custom store names is planned for the future.

---

## Usage

### Create a new AWS profile

```bash
aws-gopass-env create dev \
  --access-key AKIA... \
  --secret-access-key secret \
  --region eu-central-1
```

### List all profiles

```bash
aws-gopass-env list
```

Shows profiles stored under the `aws` namespace. The currently active one is marked with `*`.

### Show a specific profile

```bash
aws-gopass-env show dev
```

### Set environment variables from a profile

```bash
eval "$(aws-gopass-env set-env dev)"
```

This sets:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_DEFAULT_REGION`
- `AWS_PROFILE`
-  Make sure the ~/.aws/config has profile with region

You can also add this to your `~/.zshrc`:

```bash
function awsenv() {
  eval "$($HOME/path/to/aws-gopass-env set-env $1)"
}
```

So you can easily switch profiles with:

```bash
awsenv dev
```

---

### Delete a profile

```bash
aws-gopass-env delete dev
```

### Update a profile

```bash
aws-gopass-env update dev --region us-east-1
```

---

## Starship Integration

If you’re using [starship](https://starship.rs/), AWS profile information will show up automatically in your prompt when `AWS_PROFILE` or `AWS_ACCESS_KEY_ID` is set.

Customize it in your `~/.config/starship.toml`:

```toml
[aws]
format = '[$symbol($profile) ($region)]($style) '
style = "bold yellow"
```

---

## Known limitations / TODOs

- Currently only supports a `gopass` store named `aws`
- No validation to prevent overwriting existing profiles when creating
- No profile name autocomplete yet

---

## License

MIT

