install_capnp:
	git submodule update --init
	# https://capnproto.org/install.html
	cd capnproto/c++ \
	 && autoreconf -i \
	 && ./configure \
	 && make -j6 check \
	 && sudo make install

init: install_capnp
