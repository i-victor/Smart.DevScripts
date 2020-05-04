
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


project('ui-vala-glade-test', 'c', 'vala')

valadeps = [
	dependency('gtk+-3.0'),
]

executable(meson.project_name(), 'ui.vala', dependencies: valadeps)


# END
