**AUGUST - HTML Editor for LINUX/UNIX

**Version: 0.63b 2001-04-16

**Send bug reports & comments to: johanb@lls.se

**Changes in this release:

*Balloon help has replaced the status bar indicator.
*A bug in the file_save_as proc is fixed.
*Preview button removed and replaced by a Preview menu with
keyboard shortcuts.
*Experimental preview with Konqueror browser.
*Improved insertation of <A HREF tag. It one highlights an URL the code will
wrap around it.
*Insert date/time from tools menu.

**Changes in 0.62b release:

*Fixed syntax highlighting bug that caused it to miss last visible line in the
window.
*Syntax highlighting is now slightly faster.
*Added local link functionality to hyperlink button.
*Minor adjustments and additions to help file and tool tips.
*Editor line wrap setting is now properly saved.
*Syntax higlighting now works when opening files with the "-all"
command line switch.

**Changes in 0.61b release:

*The weblint syntax check window now prints weblint version number.
*Fixed bug that prevented the use of html tags with """ in the
user configurable menu.
*Changed the format of the "Augustoptions.tcl" file - you may have to
delete or edit the old file to be able to run August. See "install.txt" for
more info.

**Changes in 0.60b release:

*Syntax highlighting - basic but useful.
*August can now be configured to use the last directory
visited when opening and saving files.
*Syntax check with Weblint.
*New keyboard shortcuts: F1: Help, F2: Use template, F3: Manage user menu,
F4: Weblint check, F5: Updates syntax coloring, F6: August options.
*Set font size relative in font size/color dialog.
*August is now available as a wrapped standalone executable.
*Bugfixes.

**Changes in 0.53 release:

This is mostly a bug fix release:
*You can no longer run multiple August windows. Running multiple windows can create seriuos
confusion if the user opens the same file twice in different instances of August...;-)
If you try to start another instance of August you will just deiconify the original window.
*Undo now also works with lines that contain """, "[" or "{" characters.
Also some new features:
*Spanish and Italian special characters insertation.
*Added "*.PHP" format to the file selector in the "Open file" dialog.
*Yes, August can be used under MS Windows (Macintosh anybody?). I actually added the possibility to preview
your work under NT to this version (untested). But this is "unsupported"...;-)

**Changes in 0.52 release:

*User configurable tag menu.
*Added "*.PNG" format to the file selector in the insert image dialog.
*You can now insert the basefont tag by right clicking on the "Font size and color" button.
*Added "misc" to fontselector in options dialog.
*Bug fixes (file_close, file_save_as).

NOTE:
If you have been running August 0.42, you'd better take a look at the hidden ".August" subdirectory in your
home directory. This particular release creates a lot of buffer files that isn't erased when August is shut
down, so there may be large number of files there. Just do "rm *.bf" and they are gone. Hint: Don't do
this when August is running!

**Changes in 0.50 release:

*Improved search and replace: Search in all open files, "replace all" without prompting,
search from cursor.
*Added "Save all files" command to the "File" menu.
*August can now be started with the "-all" option, which will load all
*.html, *.htm, *.css and *.txt files in the current directory.
*It is now possible to change font and font size in the editor in the "August options"
dialog.
*It is now possible to set editor line wrap to "None", "Word" or "Char  in the "August options"
dialog.
*Improved status bar.
*New vertical button bar.
*Bug fixes.

NOTE:
If you have been running August 0.42, you'd better take a look at the hidden ".August" subdirectory in your
home directory. This particular release creates a lot of buffer files that isn't erased when August is shut
down, so there may be large number of files there. Just do "rm *.bf" and they are gone. Hint: Don't do
this when August is running!

**Changes in 0.42 release:

*The width & height of images is set automatically (with a little help from Image Magick...) in the
"Insert image" dialog.
*Removed the user adjustable limit on undo - undo is now truly unlimited!;-)
*You should now be able to run several August windows at once - I haven't tested this much so...
*Improved functionality of the "BGCOLOR" button.
*I found some really stupid bugs in the undo proc which are hopefully fixed now....
*Some additions to the helpfile.
*Other bugfixes and small improvements.

**Changes in 0.40 release:

*Unlimited undo - you can set the number of undo levels in "August options".
*New tool in the "Tools" menu - "Remove tabs". There are two possibilities:
Keep original format and remove format.
*Preview with Lynx browser. Click right mouse button on the preview button. If
it doesn't work, check the "Lynx terninal" setting in "August options". 
*Insert table rows dialog box.
*Automatic insertation of html code for country specific special characters.
So far only swedish is available. Read more about this in specchars.txt.
*Bugfixes.

**Changes in 0.37 release:

*You can now start August from the command line, ex: august "filename".
*File menu now has a "Reload file" command.
*You can insert frame tags from the Tags menu.
*Various bugfixes and small improvements. 

**Changes in 0.35 release:

*"Create table cells" dialog.
*Removed basic structure button and replaced it with
justify button. You can use the template dialog to insert
a basic structure instead.
*Improved "Font style" dialog - set font face.
*You can now set path to netscape executable in the options
dialog.
*Many bug fixes.

**Changes in 0.30 release:

* Create and use templates.
* Improved printer dialog.
* "Close all files" command on the "File" menu.
* Improved "Tags" menu.
* Change text from uppercase to lowercase and vice versa.
* Improved helpfile.
* Bugfixes.

**Changes in 0.27 release:

* Improved statusbar.
* Help file (well, at least I started working on one...;-))
* "Tools" menu with "Remove tags" feature.
* Very simple printing dialog....will be improved in future releases.
* "Go to line and "Select all" in Edit menu.
* Many bug fixes.

**Changes in 0.25 release:

* The program have changed name from xSITE to August.
* Fixed some bugs related to saving & opening files.
* I've done some work on the "Misc" menu.
* Preview with Netscape is improved - one doesn't have to start Netscape before using the
  preview button.
* Minor improvements & bugfixes.

**Main features:

* Edit multiple files.
* Search and replace.
* Support for many of the standard HTML tags.	
* Create and use templates.
* Dialog boxes for setting font size & color, inserting tables & images.

**Legal stuff...

August - HTML Editor for UNIX by Johan Bengtsson 
Copyright (C) 1999 Johan Bengtsson johanb@lls.se
Snail mail: Johan Bengtsson, Fangdammsvagen 10, 433 43 Partille, Sweden.

**NOTE: Some of the icons in the toolbar are from the free IBM/Lotus set of icons.  

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
NU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program(*); if not, write to the Free Software
Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.

(*) Read license.txt.

