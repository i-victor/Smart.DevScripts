
### Building The Project:
# cd ./this-dir
# meson builddir
# cd builddir
# ninja
###


# project definition
project('vala-project', 'c', 'vala')

# define vala packages deps ; # Sample (Direct Compile): valac --pkg gtk+-3.0 "%f"
valadeps = [
	dependency('gtk+-3.0'),
]

# define executable ; Sample: executable('prog', 'main.vala', 'class1.vala', 'class2.vala', dependencies : valadeps)
executable('hello-world', 'hello-world.vala', dependencies: valadeps)


###
# Meson Reference Samples for Vala: https://github.com/mesonbuild/meson/tree/master/test%20cases/vala
###

# END
