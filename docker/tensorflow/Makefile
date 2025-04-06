all: libtensorflow libtensorflow-avx libtensorflow-avx2 libtensorflow-vnni
# Downloads the TensorFlow source code from GitHub and extracts it to a subdirectory:
download:
	wget https://github.com/tensorflow/tensorflow/archive/v$(TF_VERSION).tar.gz
	tar -xzf v$(TF_VERSION).tar.gz
	cp Makefile *.sh tensorflow-$(TF_VERSION)
# Clang command line argument reference:
# https://clang.llvm.org/docs/ClangCommandLineReference.html#x86
libtensorflow:
	bazel build -c opt //tensorflow:libtensorflow.so --copt=-O3 --copt=-mno-avx2 --copt=-mno-avx --copt=-march=core2
	./create_archive.sh amd64 $(TF_VERSION)
libtensorflow-avx:
	bazel build -c opt //tensorflow:libtensorflow.so --copt=-O3 --copt=-mavx
	./create_archive.sh amd64-avx $(TF_VERSION)
libtensorflow-avx2:
	bazel build -c opt //tensorflow:libtensorflow.so --copt=-O3 --copt=-mavx2 --copt=-mfma --copt=-mfpmath=both --copt=-msse4.2
	./create_archive.sh amd64-avx2 $(TF_VERSION)
libtensorflow-vnni:
	bazel build -c opt //tensorflow:libtensorflow.so --copt=-O3 --copt=-mavx2 --copt=-mavxvnni --copt=-mfma --copt=-mfpmath=both --copt=-msse4.2
	./create_archive.sh amd64-vnni $(TF_VERSION)