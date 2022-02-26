#!/usr/bin/env bash
#
# Cut a podman-tui release.  Usage:
#
#   $ release.sh <version> <next-version>
#
# For example:
#
#   $ release.sh 1.2.3 1.3.0
#

VERSION="$1"
NEXT_VERSION="$2"
DATE=$(date '+%Y-%m-%d')
LAST_TAG=$(git describe --tags --abbrev=0)

write_go_version()
{
	LOCAL_VERSION="$1"
	sed -i "s/^\(.*appVersion = \"\).*/\1${LOCAL_VERSION}\"/" cmd/version.go
}

write_spec_version()
{
	LOCAL_VERSION="$1"
	LOCAL_RELEASE="1%{?dist}"
	if [[ "${LOCAL_VERSION}" == *-dev ]]; then
		LOCAL_VERSION=$(echo ${LOCAL_VERSION} | sed "s/-dev//")
		LOCAL_RELEASE="dev.1%{?dist}"
	fi
	sed -i "s/^\(Version: *\).*/\1${LOCAL_VERSION}/" podman-tui.spec.rpkg
	sed -i "s/^\(Release: *\).*/\1${LOCAL_RELEASE}/" podman-tui.spec.rpkg
}

write_spec_changelog()
{
	sed '/\*.*-dev/d' -i podman-tui.spec.rpkg
	VERSION=$1
	date=$(date "+%a %b %d %Y")
	echo "* ${date} $(git config user.name) <$(git config user.email)> ${VERSION}-1" >.changelog.txt
	if [[ "${VERSION}" != *-dev ]]; then
	   git log --no-merges --format='- %s' "${LAST_TAG}..HEAD" >>.changelog.txt
	else
	    echo "" >>.changelog.txt
	fi
	sed '/^%changelog.*/r .changelog.txt' -i podman-tui.spec.rpkg
	rm -f .changelog.txt
}

release_commit()
{
	write_go_version "${VERSION}" &&
	write_spec_version "${VERSION}" &&
	write_spec_changelog "${VERSION}" &&
	git commit -asm "Bump to v${VERSION}"
}

dev_version_commit()
{
	write_go_version "${NEXT_VERSION}-dev" &&
	write_spec_version "${NEXT_VERSION}-dev" &&
	write_spec_changelog "${NEXT_VERSION}-dev" &&
	git commit -asm "Bump to v${NEXT_VERSION}-dev"
}


git fetch origin &&
git checkout -b "bump-${VERSION}" origin/main &&
release_commit &&
git tag -s -m "version ${VERSION}" "v${VERSION}" &&
dev_version_commit
