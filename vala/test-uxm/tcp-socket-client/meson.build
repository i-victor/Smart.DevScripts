
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


project('tcp-clients', 'c', 'vala')

valadeps = [
	dependency('glib-2.0'),
	dependency('gio-2.0'),
]

executable('tcp-client-sync',  'tcp-client-sync.vala',  dependencies: valadeps)
executable('tcp-client-async', 'tcp-client-async.vala', dependencies: valadeps)

# END
