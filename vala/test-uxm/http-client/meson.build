
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


project('http-client', 'c', 'vala')

valadeps = [
	dependency('libsoup-2.4'),
]

executable(meson.project_name(), meson.project_name()+'.vala', dependencies: valadeps)


# END
