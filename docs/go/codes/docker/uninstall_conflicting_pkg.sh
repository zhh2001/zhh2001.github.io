# 只收集当前已安装的冲突包
conflicting_packages=()
for pkg in docker.io docker-compose docker-compose-v2 docker-doc docker-buildx podman-docker containerd runc; do
  if [ "$(dpkg-query -W -f='${db:Status-Status}' "$pkg" 2>/dev/null)" = 'installed' ]; then
    conflicting_packages+=("$pkg")
  fi
done

if [ "${#conflicting_packages[@]}" -gt 0 ]; then
  sudo apt remove "${conflicting_packages[@]}"
fi
