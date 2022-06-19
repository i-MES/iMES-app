#!/bin/sh

un=`uname -s`
os="unknown"
if [ $un == "Darwin" ];then
	os="macos"
elif [ $un == "Linux" ];then
	os="linux"
elif [ $un == "windows" ];then 
	os="windows"
elif [ $un == "MINGW" ];then 
	os="windows"
fi

# 1. 准备工作
#     build python in pyenv:
#     env PYTHON_CONFIGURE_OPTS="--enable-shared" pyenv install -v 3.y.z

# 2. 查找 pyenv 中有共享库的 python 版本，并检查入参合法性
vers=`find $HOME/.pyenv/versions -name "libpython3.*[so|dylib]"| \
      awk -F'versions/3.' '{print "3."$2}'| \
      awk -F'/' '{print $1}' | sort -t. -n -k2`

help() {
  echo "Usage:"
  echo "  ./change_py_pc.sh [`echo ${vers}|sed 's/ /|/g'`]"
}

if [ $# -eq 0 ];then
help
exit 1
fi

pyenv=$HOME/.pyenv/versions/${1}
if [ ! -d $pyenv ];then
help
exit 1
fi

x=${1%%.*}
yz=${1#*.}
y=${yz%.*}
z=${1##*.}
echo $x $y $z

# 3. 配置 pkg-config pc file
#    (A) 在 /usr/lib/x86_64-linux-gnu/pkgconfig/ 创建指定 pc 文件的链接
#    (B) 设置 PKG_CONFIG_PATH 环境变量 —— 存在多个 python3-embed.pc 的问题
pkgdir="."
if [ $os == "linux" ];then
pkgdir=/usr/lib/x86_64-linux-gnu/pkgconfig
elif [ $os == "macos" ];then
pkgdir=/usr/local/lib/pkgconfig
else
echo can not find pkgconfig path
exit 1
fi

sudo ln -sf $HOME/.pyenv/versions/${1}/lib/pkgconfig/python-$x.$y-embed.pc $pkgdir/python3-embed.pc
go clean --cache
echo New python pkg-config config:
echo `pkg-config --libs python3-embed`

# 4. 配置 libpythonx.6.so
#    (A) ln -s 或 cp so 到特定位置（ubutnu:/usr/lib/x86_64-linux-gnu/）—— 适合生产环境
#       find ~/.pyenv/versions -name "libpython*.so" |xargs -I{} sh -c 'sudo ln -s {} /usr/lib/x86_64-linux-gnu/'
#       ln -s $HOME/.pyenv/versions/${1}/lib/libpython* /usr/lib/x86_64-linux-gnu/
#    (B) 提示用户设置 LD_LIBRARY_PATH —— 适合 dev 环境
echo Run
cmd="LD_LIBRARY_PATH=$HOME/.pyenv/versions/${1}/lib"
echo $cmd
type xclip >/dev/null 2>&1 && echo $cmd |xclip -in -selection clipboard && echo "已拷贝到粘贴板，请复制到 wails 窗口"

# 5. Tips
echo "you may need"
echo "  pyenv shell $1"
echo "  pip install requirement package"