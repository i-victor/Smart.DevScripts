
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


# project definition
project('vala-project', 'c', 'vala')

# define vala packages deps
valadeps_gtk = [
	dependency('gtk+-3.0'),
]
valadeps_glib = [
	dependency('glib-2.0'),
]


# define executable ; Sample: executable('prog', 'main.vala', 'class1.vala', 'class2.vala', dependencies : valadeps)
executable('hello-world',  'hello-world.vala',  dependencies: valadeps_gtk)
executable('hello-world2', 'hello-world2.vala', dependencies: valadeps_gtk)
executable('hello-world3', 'hello-world3.vala', dependencies: valadeps_glib)

###
# Meson Reference Samples for Vala: https://github.com/mesonbuild/meson/tree/master/test%20cases/vala
###

# END
