package commands

import (
	"context"
	"fmt"
	"os"

	"lesiw.io/command"
	"lesiw.io/command/ctr"
	"lesiw.io/command/sys"
)

const (
	baseImage        = "alpine:3.19"
	golangVersion    = "1.24.3"
	terraformVersion = "1.8.4"
	terraformURL     = "https://releases.hashicorp.com/terraform"
	golangcilint     = "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
)

type builder struct {
	ctx context.Context
	sh  *command.Sh
	m   command.Machine
}

func (Ops) Build() error {
	ctx := context.Background()
	tag := os.Getenv("IMAGE_TAG")
	if tag == "" {
		tag = "latest"
	}
	imageTag := "kiloexabyte/runner-image:" + tag

	m := ctr.Machine(sys.Machine(), baseImage)
	defer command.Shutdown(ctx, m)

	id, err := command.Read(ctx, m, "hostname")
	if err != nil {
		return fmt.Errorf("get container id: %w", err)
	}

	b := &builder{
		ctx: ctx,
		m:   m,
		sh: command.Shell(
			m, "apk", "curl", "tar", "wget", "unzip", "go", "npm", "rm",
		),
	}

	if err := b.installPackages(); err != nil {
		return err
	}
	if err := b.installGo(); err != nil {
		return err
	}
	if err := b.installTerraform(); err != nil {
		return err
	}
	if err := b.installGoTools(); err != nil {
		return err
	}
	if err := b.installPnpm(); err != nil {
		return err
	}
	if err := b.displayVersions(); err != nil {
		return err
	}
	if err := b.cleanup(); err != nil {
		return err
	}
	if err := commitImage(ctx, id, imageTag); err != nil {
		return err
	}

	return nil
}

func (b *builder) installPackages() error {
	fmt.Println("Installing system packages...")
	return b.sh.Exec(b.ctx, "apk", "add", "--no-cache",
		"bash", "curl", "git", "docker-cli", "ca-certificates",
		"nodejs", "npm", "tar", "zip", "aws-cli", "gcc", "musl-dev",
	)
}

func (b *builder) installGo() error {
	fmt.Println("Installing Go...")
	tarball := fmt.Sprintf("go%s.linux-amd64.tar.gz", golangVersion)
	url := fmt.Sprintf("https://golang.org/dl/%s", tarball)
	if err := b.sh.Exec(b.ctx, "curl", "-LO", url); err != nil {
		return err
	}
	if err := b.sh.Exec(b.ctx,
		"tar", "-C", "/usr/local", "-xzf", tarball); err != nil {
		return err
	}
	return b.sh.Exec(b.ctx, "rm", tarball)
}

func (b *builder) installTerraform() error {
	fmt.Println("Installing Terraform...")
	zip := fmt.Sprintf("terraform_%s_linux_amd64.zip", terraformVersion)
	url := fmt.Sprintf("%s/%s/%s", terraformURL, terraformVersion, zip)
	if err := b.sh.Exec(b.ctx, "wget", url); err != nil {
		return err
	}
	if err := b.sh.Exec(b.ctx,
		"unzip", zip, "-d", "/usr/local/bin"); err != nil {
		return err
	}
	return b.sh.Exec(b.ctx, "rm", zip)
}

func (b *builder) installGoTools() error {
	fmt.Println("Installing Go tools...")
	if err := b.sh.Exec(b.ctx,
		"go", "install", "lesiw.io/op@latest"); err != nil {
		return err
	}
	return b.sh.Exec(b.ctx, "go", "install", golangcilint+"@v2.1.6")
}

func (b *builder) installPnpm() error {
	fmt.Println("Installing pnpm...")
	return b.sh.Exec(b.ctx, "npm", "install", "-g", "pnpm")
}

func (b *builder) displayVersions() error {
	fmt.Println("Displaying versions...")
	if err := b.sh.Exec(b.ctx, "go", "version"); err != nil {
		return err
	}
	if err := command.Do(b.ctx, b.m, "node", "-v"); err != nil {
		return err
	}
	return command.Do(b.ctx, b.m, "pnpm", "-v")
}

func (b *builder) cleanup() error {
	fmt.Println("Cleaning up caches...")
	if err := b.sh.Exec(b.ctx,
		"go", "clean", "-cache", "-modcache"); err != nil {
		return err
	}
	if err := b.sh.Exec(b.ctx,
		"npm", "cache", "clean", "--force"); err != nil {
		return err
	}
	return b.sh.Exec(b.ctx, "rm", "-rf", "/var/cache/apk/*")
}

func commitImage(ctx context.Context, id, imageTag string) error {
	fmt.Println("Committing image...")
	ctl := sys.Machine()
	if err := command.Do(ctx, ctl,
		"docker", "container", "commit",
		"--change", "ENV GOLANG_VERSION="+golangVersion,
		"--change", "ENV GOROOT=/usr/local/go",
		"--change", "ENV GOPATH=/go",
		"--change", "ENV PATH=/usr/local/go/bin:/go/bin:$PATH",
		id, imageTag,
	); err != nil {
		return err
	}
	fmt.Printf("Built image: %s\n", imageTag)
	sh := command.Shell(ctl, "docker")
	return sh.Exec(ctx, "docker", "images", imageTag,
		"--format", "Image Size: {{.Size}}")
}
