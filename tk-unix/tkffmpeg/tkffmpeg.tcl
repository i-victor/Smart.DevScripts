#!/bin/sh
#The next line executes wish - wherever it is \
exec /usr/local/bin/wish8.5 "$0" "$@"
#
# tkffmpeg - frontend for ffmpeg
# tkffmpeg comes with ABSOLUTELY NO WARRANTY; for details read license.
# This is free software, and you are welcome to redistribute it
# under certain conditions.
#  
package require Tk
package require Ttk
#package require img::png
#
puts "tkffmpeg - frontend for ffmpeg\n"
# Procedures
proc inFile {} {
	global dirIn inpF
	set types {
		{{All Files} * }
	} 
	if {$dirIn==""} {
		set dirIn [string range $inpF 0 [string last "/" $inpF]]
	}
	set filename [tk_getOpenFile -filetypes $types -title "Select Input File"  -initialdir $dirIn]
	set dirIn ""
	.fiLes.fiLes1.textFile delete 0 end
	.fiLes.fiLes1.textFile insert end $filename
	.fiLes.fiLes1.textFile xview [string lengt $filename]
	.log.msg insert end [getinfo $filename]
	return
}
proc sFile {} {
	global dirOut
	set filename [tk_getSaveFile -title "Select Output File" -initialdir $dirOut]
	.fiLes.fiLes2.textFile delete 0 end 
	.fiLes.fiLes2.textFile insert end $filename
	.fiLes.fiLes2.textFile xview [string lengt $filename]
	return
}
proc TempFolder {} {
	global tempdir
	if {![file exists $tempdir]} {
		set tempdir "~"
	}
	set dir [tk_chooseDirectory -initialdir $tempdir -title "Choose a directory"]
	.fiLes.vframe.temp delete 0 end
	.fiLes.vframe.temp insert end  $dir
	.fiLes.vframe.temp xview [string lengt $dir]
	return 
}
proc MultiFolder {} {
	global multiple
	if {![file exists $multiple]} {
		set multiple "~"
	}
	set dir [tk_chooseDirectory -initialdir $multiple -title "Choose a directory"]
	.fiLes.packed.temp delete 0 end
	.fiLes.packed.temp insert end  $dir
	.fiLes.packed.temp xview [string lengt $dir]
	return 
}
proc log {} {
	global input ln ok_check openFu outF error_run programm end time_run lang_var
	if { [gets $input line] > 0 } {
		.log.msg insert end $line\n
		.log.msg see end
		set ln $line
		update
	} elseif { [eof $input] } {
		catch { close $input }
		set newline [split $ln :]
		if {[lindex $newline 0]=="video"} {
			set ok_check 1
		} 
		if {$ok_check==1} {
			.log.msg insert end "$lang_var(lang_ok)\n" blue_color
			set end 1
		} elseif {$ok_check==0 && $error_run==0} {
			.log.msg insert end "$lang_var(error_log)\n" red_color	
			set end 1
		}
		if {$openFu==1 && [file exists $outF] && $error_run!=1} {
			if {$programm=="" || ![file exist "/usr/local/bin/$programm"]} {
				.log.msg insert end "$lang_var(error_programm)\n" red_color
			} elseif  {$ok_check==1} {
				open "|$programm \"$outF\""
				.log.msg insert end "$lang_var(lang_open_file)\n" blue_color
			}
		}
		#.but.go configure -state normal
		set second [expr {[clock seconds]-$time_run}]
		if {$second>60} {
			set minute [expr {$second/60}]
			set secondus [expr {$second - $minute*60}]
			if {$minute>60} {
				set hour [expr {$minute/60}]
				set minutes [expr {$minute-$hour*60}]
				set time_text "$hour h. $minutes min. $secondus sec."
			} else {
				set time_text "$minute min. $secondus sec."
			}
		} else {
			set time_text "$second sec."
		}
		.log.msg insert end "$lang_var(lang_time) $time_text\n"
		.log.msg see end
		Progres "stop"
	}
}
proc go {inpF outF vBit fps chdein x y vcodec acodec freq aBit st more_opt} {
	global input ln pi openFu error_run programm ok_check lunch time_run lang_var acopy vcopy
	set error_run 0
	set lunch "go"
	set time_run [clock seconds]
	Progres "start"
	#.log.msg delete 1.0 end
	#.but.burn configure -state disabled
	.but.go configure -state disabled
	if {$inpF=="" && $outF==""} {Progres "stop";return "$lang_var(error_file)\n"}
	if {$inpF==""} {Progres "stop";return "$lang_var(error_input_file)\n"}
	if {![file exists "$inpF" ]} {Progres "stop";return "$lang_var(error_input_file_ex)\n"}
	if {$outF==""} {Progres "stop";return "$lang_var(error_output_file)\n"}
	if {$chdein==1} {
		set chde {-deinterlace}
	} else {
		set chde {}
	}
	if {$st==1} {
		set chst {-ac 2}
	} else {
		set chst {-ac 1}
	}
	if {$vBit=="" || $vBit==0} {Progres "stop";return "$lang_var(error_vbitrate)\n"}
	if {$x=="" || $x==0} {Progres "stop";return "$lang_var(error_x)\n"}
	if {$y=="" || $y==0} {Progres "stop";return "$lang_var(error_y)\n"}
	if {$fps=="" || $fps==0} {Progres "stop";return "$lang_var(error_fps)\n"}
	if {$vcodec=="" || $vcodec==0} {Progres "stop";return "$lang_var(error_vcodec)\n"}
	if {$acodec=="" || $acodec==0} {Progres "stop";return "$lang_var(error_acodec)\n"}
	if {$freq=="" || $freq==0} {Progres "stop";return "$lang_var(error_freq)\n"}
	if {$aBit=="" || $aBit==0} {Progres "stop";return "$lang_var(error_abitrate)\n"}
	if {[file exists "$outF" ]} {
		set chfile "-y"
	} else {
		set chfile ""
	}
	if {$acodec=="mp3"} {
		set acodec "libmp3lame"
	}
	if {$vcodec=="libx264" && $more_opt==""} {
		set more_opt "-vpre libx264-default"
		.addit.more insert end $more_opt
	}
	set str "ffmpeg -i \"$inpF\" -b $vBit\k -r $fps $chde -s $x\x$y -vcodec $vcodec $more_opt $chst -acodec $acodec -ar $freq -ab $aBit\k $chfile \"$outF\""
	if {$vcopy==1} {
            set str "ffmpeg -i \"$inpF\" -vcodec copy $chst -acodec $acodec $more_opt -ar $freq -ab $aBit\k $chfile \"$outF\""
        } elseif {$acopy==1} {
            set str "ffmpeg -i \"$inpF\" -b $vBit\k -r $fps $chde -s $x\x$y -vcodec $vcodec $more_opt -acodec copy $chfile \"$outF\""
        } elseif {$acopy==1 && $vcopy==1} {
            set str "ffmpeg -i \"$inpF\" -vcodec copy -acodec copy $more_opt $chfile \"$outF\""
        }
        set ok_check 0	
	if {[catch {open "|$str |& cat"} input]} {
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
		} else {
			set pi [lindex [pid $input] 0]
			fileevent $input readable log
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
	}
	catch {fconfigure $input -blocking 0}
	.log.msg see end
}
proc kill {} {
	global pi error_run lang_var lunch
	set ex 0
	if {$pi!=""} {
		set fill [open "|kill -s kill $pi"]
	}
	set error_run 1
	catch {destroy .progres}
	.log.msg insert end "$lang_var(error_user)\n" red_color
	.log.msg see end
	if {$lunch==""} {destroy .}
	if {$lunch!="go"} {
		.but.burn configure -state active
	}
	.but.go configure -state active
}
proc save_options {} {
	global env programm inpF outF vBit fps dIn vcodec X Y acodec freq aBit stu openFu language version more_opt
	global userO rundir tempdir multiple
	if { ! [file exists "$env(HOME)/.Tkffmpeg"]} {exec touch $env(HOME)/.Tkffmpeg}
	set config_file [open $env(HOME)/.Tkffmpeg w]
	puts $config_file "programm = $programm"
	set dirIn [string range $inpF 0 [string last "/" $inpF]]
	puts $config_file "dirIn = $dirIn"
	set dirOut [string range $outF 0 [string last "/" $outF]]
	puts $config_file "dirOut = $dirOut"
	puts $config_file "language = $language"
	puts $config_file "userO = $userO"
	puts $config_file "vBit = $vBit"
	puts $config_file "fps = $fps"
	puts $config_file "dIn = $dIn"
	puts $config_file "vcodec = $vcodec"
	puts $config_file "X = $X"
	puts $config_file "Y = $Y"
	puts $config_file "acodec = $acodec"
	puts $config_file "freq = $freq"
	puts $config_file "aBit = $aBit"
	puts $config_file "stu = $stu"
	puts $config_file "openFu = $openFu"
	puts $config_file "tempdir = $tempdir"
	puts $config_file "more_opt = $more_opt"
	puts $config_file "multiple = $multiple"
	close $config_file
	file delete -force "$rundir"
	exit
}
proc adv_opt_window {} {
	global language programm userO user lang_var
	if { ! [winfo exists .adv_opt_window] } {
		toplevel .adv_opt_window
		wm title .adv_opt_window $lang_var(lang_adv_opt)
		wm resizable .adv_opt_window false false
		wm transient .adv_opt_window .
		set lang_use [lang_list]
		frame .adv_opt_window.1 -bg #b4d9cf
		label .adv_opt_window.1.llang -text $lang_var(lang_change) -takefocus 0 -font {-size 10} -bg #b4d9cf
		ttk::combobox .adv_opt_window.1.lang -values $lang_use -width 17 -textvariable language -state readonly
		.adv_opt_window.1.lang current [lsearch -exact $lang_use $language]
		frame .adv_opt_window.2 -bg #b4d9cf
		label .adv_opt_window.2.lprog -text $lang_var(lang_programm) -takefocus 0 -font {-size 10} -bg #b4d9cf
		ttk::entry .adv_opt_window.2.prog -width 10 -textvariable programm
		frame .adv_opt_window.3 -bg #b4d9cf
		checkbutton .adv_opt_window.3.userO -variable userO -text $lang_var(lang_userO) -bg #b4d9cf
		frame .adv_opt_window.4 -bg #b4d9cf
		ttk::button .adv_opt_window.4.ok -command {destroy .adv_opt_window; array set lang_var [load_lang $language];reload_primary;} -text $lang_var(lang_save_but) -width 10
		pack .adv_opt_window.1 -fill x
		pack .adv_opt_window.1.llang -side left -padx 5 -pady 2
		pack .adv_opt_window.1.lang -side left -padx {0 5}
		pack .adv_opt_window.2 -fill x
		pack .adv_opt_window.2.lprog -side left -pady 2 -padx 5
		pack .adv_opt_window.2.prog -side left -padx {0 5}
		pack .adv_opt_window.3 -fill x
		pack .adv_opt_window.3.userO -side left -pady 2 -padx 5
		pack .adv_opt_window.4 -fill x
		pack .adv_opt_window.4.ok
	}
}
proc lang_list {} {
	global path_proc
	set defFile "$path_proc/language"
	if { [file exists $defFile] } {
            foreach value [glob -nocomplain -types f -directory $defFile *] {
                lappend lang_use [file rootname [file tail $value]]
            }
	}
	return $lang_use
}
proc dvdcheck {check} {
	global lang_var
	if {$check==1} {
		.fiLes.fiLes2.name configure -state disabled
		.fiLes.fiLes2.textFile configure -state disabled
		.fiLes.fiLes2.saveFile configure -state disabled
		.fiLes.packed.dvdvideo configure -state disabled
		.fiLes.packed.ltemp configure -state disabled
		.option.optionV.v1.4.label configure -state disabled
		.option.optionV.v1.4.lst configure -state disabled
		.option.optionA.1.lstac configure -state disabled
		.option.optionA.1.ac configure -state disabled
		.but.chop configure -state disabled
		#.option.optionV.v1.1.bitrate configure -state disabled
		#.option.optionV.v1.1.sbitrate configure -state disabled
		#.option.optionV.v1.2.sfps configure -state disabled
		#.option.optionV.v1.2.fps configure -state disabled
		#.option.optionA.2.fr configure -state disabled
		#.option.optionA.2.lstfr configure -state disabled
		#.option.optionA.3.bit configure -state disabled
		#.option.optionA.3.lstbit configure -state disabled
		.but.go configure -command {.log.msg insert end [dvdgo $inpF $vBit $fps $dIn $freq $aBit $stu $aspect] red_color}
		ttk::button .but.burn -command burn -text $lang_var(.but.burn) -width 7 -state disabled
		pack .but.burn -side right -padx {0 5}
		.fiLes.vframe.ltemp configure -state normal
		.fiLes.vframe.temp configure -state normal
		.fiLes.vframe.tempfolder configure -state normal
		destroy .option.optionV.size.2
		destroy .option.optionV.size.3
		#
		ttk::frame .option.optionV.size.2  
		pack .option.optionV.size.2 -pady 5 -fill x -padx 10
		ttk::label .option.optionV.size.2.x -text $lang_var(.option.optionV.size.2.x) -takefocus 0 -font {-size 10}
		pack .option.optionV.size.2.x -fill x
		ttk::combobox .option.optionV.size.2.aspect -values {"4:3" "16:9"} -width 5 -textvariable aspect
		.option.optionV.size.2.aspect current 0		
		pack .option.optionV.size.2.aspect -fill x
	} else {
		.fiLes.fiLes2.name configure -state normal
		.fiLes.fiLes2.textFile configure -state normal
		.fiLes.fiLes2.saveFile configure -state normal
		.fiLes.packed.dvdvideo configure -state normal
		.fiLes.packed.ltemp configure -state normal
		.option.optionV.v1.4.label configure -state normal
		.option.optionV.v1.4.lst configure -state normal
		.option.optionA.1.lstac configure -state normal
		.option.optionA.1.ac configure -state normal
		.but.chop configure -state normal
		#.option.optionV.v1.1.bitrate configure -state normal
		#.option.optionV.v1.1.sbitrate configure -state normal
		#.option.optionV.v1.2.sfps configure -state normal
		#.option.optionV.v1.2.fps configure -state normal
		.but.go configure -command {.log.msg insert end [go $inpF $outF $vBit $fps $dIn $X $Y $vcodec $acodec $freq $aBit $stu $more_opt] red_color}
		destroy .but.burn
		.fiLes.vframe.ltemp configure -state disabled
		.fiLes.vframe.temp configure -state disabled
		.fiLes.vframe.tempfolder configure -state disabled
		#.option.optionA.2.fr configure -state normal
		#.option.optionA.2.lstfr configure -state normal
		#.option.optionA.3.bit configure -state normal
		#.option.optionA.3.lstbit configure -state normal
		#
		destroy .option.optionV.size.2
		ttk::frame .option.optionV.size.2  
		pack .option.optionV.size.2 -pady 5 -fill x -padx 10
		spinbox .option.optionV.size.2.sx -from 10 -to 1024 -increment 2 -width 5 -textvariable X
		ttk::label .option.optionV.size.2.x1 -text "X" -takefocus 0 -font {-size 10}
		pack .option.optionV.size.2.sx -side left
		pack .option.optionV.size.2.x1 -side left -padx {5 0}
		ttk::frame .option.optionV.size.3  
		pack .option.optionV.size.3 -pady 5 -fill x -padx 10
		spinbox .option.optionV.size.3.sy -from 10 -to 1024 -increment 2 -width 5 -textvariable Y
		ttk::label .option.optionV.size.3.y -text "Y" -takefocus 0 -font {-size 10}
		pack .option.optionV.size.3.sy -side left
		pack .option.optionV.size.3.y -side left -padx {5 0}
	}
}
#
proc multicheck {check} {
	if {$check==1} {
		.fiLes.fiLes1.name configure -state disabled
		.fiLes.fiLes1.textFile configure -state disabled
		.fiLes.fiLes1.inputFile configure -state disabled
		.fiLes.fiLes2.name configure -state disabled
		.fiLes.fiLes2.textFile configure -state disabled
		.fiLes.fiLes2.saveFile configure -state disabled
		.fiLes.vframe.dvdvideo configure -state disabled
		.fiLes.vframe.ltemp configure -state disabled
		.but.go configure -command {MultiGo}
		#ttk::button .but.burn -command burn -text $lang_burn -width 7 -state disabled
		#pack .but.burn -side right -padx {0 5}
		.fiLes.packed.ltemp configure -state normal
		.fiLes.packed.temp configure -state normal
		.fiLes.packed.tempfolder configure -state normal
	} else {
		.fiLes.fiLes1.name configure -state normal
		.fiLes.fiLes1.textFile configure -state normal
		.fiLes.fiLes1.inputFile configure -state normal
		.fiLes.fiLes2.name configure -state normal
		.fiLes.fiLes2.textFile configure -state normal
		.fiLes.fiLes2.saveFile configure -state normal
		.fiLes.vframe.dvdvideo configure -state normal
		.fiLes.vframe.ltemp configure -state normal
		.but.go configure -command {.log.msg insert end [go $inpF $outF $vBit $fps $dIn $X $Y $vcodec $acodec $freq $aBit $stu $more_opt] red_color}
		.fiLes.packed.ltemp configure -state disabled
		.fiLes.packed.temp configure -state disabled
		.fiLes.packed.tempfolder configure -state disabled
	}
}
proc vcopy {check} {
    if {$check==1} {
        .option.optionV.v1.1.sbitrate configure -state disabled
        .option.optionV.v1.1.bitrate configure -state disabled
        .option.optionV.v1.2.sfps configure -state disabled
        .option.optionV.v1.2.fps configure -state disabled
        .option.optionV.v1.3.chdein configure -state disabled
        .option.optionV.v1.4.lst configure -state disabled
        .option.optionV.v1.4.label configure -state disabled
        .option.optionV.size.2.sx configure -state disabled
        .option.optionV.size.2.x1 configure -state disabled
        .option.optionV.size.3.sy configure -state disabled
        .option.optionV.size.3.y configure -state disabled
    } else {
        .option.optionV.v1.1.sbitrate configure -state normal
        .option.optionV.v1.1.bitrate configure -state normal
        .option.optionV.v1.2.sfps configure -state normal
        .option.optionV.v1.2.fps configure -state normal
        .option.optionV.v1.3.chdein configure -state normal
        .option.optionV.v1.4.lst configure -state normal
        .option.optionV.v1.4.label configure -state normal
        .option.optionV.size.2.sx configure -state normal
        .option.optionV.size.2.x1 configure -state normal
        .option.optionV.size.3.sy configure -state normal
        .option.optionV.size.3.y configure -state normal
    }
}
proc acopy {check} {
    if {$check==1} {
        .option.optionA.1.ac configure -state disabled
        .option.optionA.1.lstac configure -state disabled
        .option.optionA.2.fr configure -state disabled
        .option.optionA.2.lstfr configure -state disabled
        .option.optionA.3.bit configure -state disabled
        .option.optionA.3.lstbit configure -state disabled
        .option.optionA.4.chst configure -state disabled
    } else {
        .option.optionA.1.ac configure -state normal
        .option.optionA.1.lstac configure -state normal
        .option.optionA.2.fr configure -state normal
        .option.optionA.2.lstfr configure -state normal
        .option.optionA.3.bit configure -state normal
        .option.optionA.3.lstbit configure -state normal
        .option.optionA.4.chst configure -state normal
    }
}
#
proc burn {} {
	global env rundir burn
	Progres "start"
        set burn 1
	.but.burn configure -state disabled
	.but.go configure -state disabled
	set fn [open "|du -b $rundir"]
	while {![eof $fn]} {
		if { [gets $fn line] > 0 } {
			if {[lindex $line end]==$rundir} {
				set sizedir [lindex $line 0]
			}
		}
		update
	}
	if {$sizedir<4700307456} {
		set str "growisofs -Z /dev/dvd -dvd-video $rundir"
		.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
		set fn [open "|$str |& cat"]
		fconfigure $fn -blocking 0
		while {![eof $fn]} {
			if { [gets $fn line] > 0 } {
				.log.msg insert end $line\n
				.log.msg see end
			}
			update
		}
		close $fn
		.log.msg insert end "$lang_var(lang_ok)\n" blue_color
	} else {
		.log.msg insert end "$lang_var(error_size)\n" red_color
	}
	.but.burn configure -state normal
	.but.go configure -state normal
	.log.msg see end
	Progres "stop"
}
proc logdvd {} {
	global input ln ok_check error_run lang_var
	if { [gets $input line] > 0 } {
		.log.msg insert end $line\n
		.log.msg see end
		set ln $line
		update
	} elseif { [eof $input] } {
		catch { close $input }
		if {[lindex $ln 3]=="VOBUS"} {
			set ok_check 1
		}
		if {$ok_check==1} {
			.log.msg insert end "$lang_var(lang_ok)\n" blue_color
			set str "dvdauthor -o . -T"
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
			set fn [open "|$str |& cat"]
			fconfigure $fn -blocking 0
			while {![eof $fn]} {
				if { [gets $fn line] > 0 } {
					.log.msg insert end $line\n
					.log.msg see end
				}
				update
			}
			close $fn
			.log.msg insert end "$lang_var(lang_ok)\n" blue_color
		} elseif {$ok_check==0 && $error_run==0} {
			.log.msg insert end "$lang_var(error_log)\n" red_color	
		}
		.but.burn configure -state active
		.log.msg see end
		Progres "stop"
	}
}

#Convert the DVD-video
proc dvdgo {inpF vBit fps chdein freq aBit st aspect} {
#dvdgo $inpF $outF $vBit $fps $dIn $freq $aBit $stu
	global pi tempdir rundir error_run input ln ok_check lunch lang_var
	set error_run 0
	set lunch "dvdgo"
	Progres "start"
	.log.msg delete 1.0 end
	.but.burn configure -state disabled
	.but.go configure -state disabled
	if {$inpF==""} {Progres "stop";return "$lang_var(error_input_file)\n"}
	if {$vBit=="" || $vBit==0} {Progres "stop";return "$lang_var(error_vbitrate)\n"}
	if {![file exists "$inpF" ]} {Progres "stop";return "$lang_var(error_input_file_ex)\n"}
	if {$chdein==1} {
		set chde {-deinterlace}
	} else {
		set chde {}
	}
	if {$st==1} {
		set chst {-ac 2}
	} else {
		set chst {-ac 1}
	}
	if {$vBit=="" || $vBit==0} {Progres "stop";return "$lang_var(error_vbitrate)\n"}
	if {$fps=="" || $fps==0} {Progres "stop";return "$lang_var(error_fps)\n"}
	if {$freq=="" || $freq==0} {Progres "stop";return "$lang_var(error_freq)\n"}
	if {$aBit=="" || $aBit==0} {Progres "stop";return "$lang_var(error_abitrate)\n"}
	set rundir "$tempdir/tkffmpeg"
	file delete -force "$rundir"
	file mkdir "$rundir"
	cd "$rundir"
	set str "ffmpeg -i \"$inpF\" -b $vBit\k -r $fps $chde -s 720x576 -aspect $aspect -vcodec mpeg2video $chst -acodec libmp3lame -ar $freq -ab $aBit\k -f dvd - | dvdauthor -o . -t -"
	#set str "ffmpeg -i \"$inpF\" -b $vBit\k -r $fps $chde -s $x\x$y -vcodec $vcodec $chst -acodec $acodec -ar $freq -ab $aBit\k $chfile \"$outF\""
	set ok_check 0
	if {[catch {open "|$str |& cat"} input]} {
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
		} else {
			set pi [lindex [pid $input] 1]
			fileevent $input readable logdvd
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
	}
	catch {fconfigure $input -blocking 0}
	.log.msg see end
}
# Get Info
proc getinfo {inpF} {
	global input ok_check lang_var
	.log.msg delete 1.0 end
	if {$inpF==""} {return "$lang_var(error_input_file)\n"}
	if {![file exists "$inpF" ]} {return "$lang_var(error_input_file_ex)\n"}
	set str "ffmpeg -i \"$inpF\""
	set ok_check 0
	if {[catch {open "|$str |& cat"} input]} {
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
		} else {
			set pi [lindex [pid $input] 1]
			fileevent $input readable loginfo
			.log.msg insert end "$lang_var(lang_exec) $str\n\n" blue_color
	}
	catch {fconfigure $input -blocking 0}
	.log.msg see end
}
proc loginfo {} {
	global input ln ok_check error_run lang_var
	if { [gets $input line] > 0 } {
		.log.msg insert end $line\n
		.log.msg see end
		set ln $line
		if {[lindex $ln 0]=="Duration:"} {
			set ok_check 1
			set templn [lindex [split $ln ,] end]
			if {[lindex $templn 0]=="bitrate:"} {
				.fiLes.info.1.2.vbitX configure -text "[lindex $templn 1] kb/s"
			}
		}
		if {[lindex $ln 0]=="Stream" && [string first "#0.0" [lindex $ln 1]]!=-1} {
			set templn [split $ln ,]
			.fiLes.info.1.1.vcodecX configure -text [lindex [lindex $templn 0] end]
			.fiLes.info.1.3.sizeX configure -text [lindex [lindex $templn 2] 0]
			.fiLes.info.1.4.aspectX configure -text [string trimright [lindex [lindex $templn 2] end] "]"]
		}
		if {[lindex $ln 0]=="Stream" && [string first "#0.1" [lindex $ln 1]]!=-1} {
			set templn [split $ln ,]
			.fiLes.info.2.1.acodecX configure -text [lindex [lindex $templn 0] end]
			.fiLes.info.2.2.abitX configure -text [lindex $templn end]
			.fiLes.info.2.3.frequencyX configure -text [lindex $templn 1]
			.fiLes.info.2.chst configure -text [lindex $templn 2]
		}
		update
	} elseif { [eof $input] } {
		catch { close $input } 
		if {$ok_check==1} {
			.log.msg insert end "$lang_var(lang_ok)\n" blue_color
		} elseif {$ok_check==0 && $error_run==0} {
			.log.msg insert end "$lang_var(error_log)\n" red_color	
		}
		.log.msg see end
	}
}
#
proc Progres {value} {
	global error_run lang_var progres burn
	if {$value=="start"} {
		set w .progres
		toplevel $w
		wm title $w "Operation in progress"
		wm iconname $w "tkffmpeg"
		wm resizable $w false false
		wm geometry $w 300x70
		wm transient $w .
		pack [ttk::frame $w.prog] -fill x -pady {0 5} -padx 10
		ttk::progressbar $w.prog.p2 -mode indeterminate ;#-variable i
		pack $w.prog.p2 -padx 20 -fill x -expand 1 -pady 10
		ttk::button $w.prog.but -command kill -text $lang_var(lang_abort) -width 10 
		pack $w.prog.but -anchor center
		$w.prog.p2 start
                set progres 1
	} else {
		if {$error_run!=1} {
			destroy .progres
		}
		.but.go configure -state active
                set progres 0
                set burn 0
	}
}
proc MultiGo {} {
	global multiple vBit fps dIn X Y vcodec acodec freq aBit stu end error_run lang_var f more_opt
	set i 0
	set convertdir "$multiple/convert/"
	file mkdir $convertdir
	set times_run [clock seconds]
	foreach f [lsort -dictionary [glob -nocomplain -types f -dir "$multiple" *]] {
		set end 0
		set error_run 0
		set outfile "$convertdir[file tail $f]"
		#.log.msg insert end $f\n
			.log.msg insert end [go $f $outfile $vBit $fps $dIn $X $Y $vcodec $acodec $freq $aBit $stu $more_opt] red_color
			while {$i!=1} {
				if {$end==1 || $error_run==1} {
					break					
				}
				update
				after 100
			}
		if {$error_run==1} {
			break
		}
	}
	set second [expr {[clock seconds]-$times_run}]
	if {$second>60} {
		set minute [expr {$second/60}]
		set secondus [expr {$second - $minute*60}]
		if {$minute>60} {
			set hour [expr {$minute/60}]
			set minutes [expr {$minute-$hour*60}]
			set time_text "$hour h. $minutes min. $secondus sec."
		} else {
			set time_text "$minute min. $secondus sec."
		}
	} else {
		set time_text "$second sec."
	}
	.log.msg insert end "$lang_var(lang_all_time) $time_text\n"
	.log.msg see end
}
proc load_lang {lang} {
	global path_proc
	set defFile "$path_proc/language/$lang.lng"
	if { [file exists $defFile] } {
		set config_file [open $defFile r]
		while { ! [eof $config_file] } {
			gets $config_file string_temp
			lappend lang_var_list [lindex $string_temp 0] [lrange $string_temp 2 end]
		}
		close $config_file
	}
	return $lang_var_list
}
proc all_widgets {{root .}} {
	lappend all $root
	foreach child [winfo children $root] {
		if {$child ne ".ic.ic"} {
			lappend all {*}[all_widgets $child] 
		}
	}
	return $all
}
proc search_widget_text {} {
	foreach widget [all_widgets] {
		foreach wi_text [$widget configure] {
			if {[lindex $wi_text 0]=="-text"} {
				lappend widget_text $widget
			}
		}
	}
	return $widget_text
}
proc reload_primary {} {
	global lang_var
	foreach widget [search_widget_text] {
		foreach lang [array names lang_var] {
			if {$widget==$lang} {
				if {[$widget configure -text]!=$lang_var($lang)} {
					$widget configure -text $lang_var($lang)
					update;
				}
			}
		}
	}
}
proc file_exist {name} {
	if { [file exists "./$name"] } {
		return "./$name"
	} elseif {[file exists "/opt/tk-unix/tkffmpeg/$name"]} {
		return "/opt/tk-unix/tkffmpeg/$name"
	}
}
#Main programm
# Set variable
set dirIn ""
set dirOut ""
set inpF ""
set outF ""
set pi ""
set user 0
set dvd 0
set error_run 0
set tempdir "/home/$env(USER)/tmp"
set multiple ""
set rundir ""
set language ""
set lunch ""
set progres 0
set f ""
set burn 0
set vcopy 0
set acopy 0
set more_opt ""
set path_proc [file dirname [file normalize [info script] ]]
catch {tk::IconList}
set ::tk::dialog::file::showHiddenBtn 1
set ::tk::dialog::file::showHiddenVar 0

## unixman
#image create photo ico -file "$path_proc/tkffmpeg16.png";
#image create photo icoBig -file "$path_proc/tkffmpeg.png";
##

#Bind key
bind all <Control-q> {kill} 
#Load default.ini
set defFile "$path_proc/default.ini"
if { [file exists $defFile] } {
	set config_file [open $defFile r]
	while { ! [eof $config_file] } {
		gets $config_file string_temp
		set [lindex $string_temp 0] [lrange $string_temp 2 end]
	}
	close $config_file
}
#Load saved options
if { [file exists "$env(HOME)/.Tkffmpeg"] } {
	set config_file [open $env(HOME)/.Tkffmpeg r]
	while { ! [eof $config_file] } {
		gets $config_file string_temp
		if {[lindex $string_temp end]!="=" && [lindex $string_temp end]!="-1"} {set [lindex $string_temp 0] [lrange $string_temp 2 end]}

		if {[lindex $string_temp 0]=="userO" && [lindex $string_temp end]==0} {break}
	}
	close $config_file
}
#Loading linguistic variables
if {$language!=""} {#;} else {
	set locale [split $env(LANG) .]
	switch [lindex $locale 0] {
		default {set language "english"}
	}
}
array set lang_var [load_lang $language]
#Window settings
wm title . "$lang_var(.ffmpeg) $version"
wm iconname . "tkffmpeg"
wm geometry . 590x599
wm resizable . false false
wm protocol . WM_DELETE_WINDOW { save_options }

## unixman
#wm iconphoto . -default ico icoBig;
##

ttk::style theme use alt

## unixman
#if {![catch {package require tktray} ec]} {
#	source "$path_proc/tray.tcl";    
#}
##

#GUI
ttk::label .ffmpeg -text $lang_var(.ffmpeg) -justify center -takefocus 0 -font {-size 22} -anchor center
pack .ffmpeg -fill x -ipady 15
# Section Files
ttk::labelframe .fiLes -text $lang_var(.fiLes)
ttk::frame .fiLes.fiLes1
ttk::frame .fiLes.fiLes2
# Section Info Files
frame .fiLes.info 
frame .fiLes.info.1
frame .fiLes.info.1.1
frame .fiLes.info.1.2
frame .fiLes.info.1.3
frame .fiLes.info.1.4
ttk::label .fiLes.info.1.1.vcodec -text "$lang_var(.fiLes.info.1.1.vcodec) -" -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.1.1.vcodecX -text 0 -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.1.2.vbit -text "$lang_var(.fiLes.info.1.2.vbit) -" -takefocus 0 -font {-size 10}
ttk::label .fiLes.info.1.2.vbitX -text "0 kb/s" -takefocus 0 -font {-size 10}  
ttk::label .fiLes.info.1.3.size -text "$lang_var(.fiLes.info.1.3.size) -" -takefocus 0 -font {-size 10}
ttk::label .fiLes.info.1.3.sizeX -text "0" -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.1.4.aspect -text "$lang_var(.fiLes.info.1.4.aspect) -" -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.1.4.aspectX -text "0" -takefocus 0 -font {-size 10} 
frame .fiLes.info.2
frame .fiLes.info.2.1
frame .fiLes.info.2.2
frame .fiLes.info.2.3
ttk::label .fiLes.info.2.1.acodec -text "$lang_var(.fiLes.info.2.1.acodec) -" -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.2.1.acodecX -text "0" -takefocus 0 -font {-size 10} -width 10 
ttk::label .fiLes.info.2.2.abit -text "$lang_var(.fiLes.info.2.2.abit) -" -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.2.2.abitX -text "0" -takefocus 0 -font {-size 10} -width 10 
ttk::label .fiLes.info.2.3.frequency -text "$lang_var(.fiLes.info.2.3.frequency) -" -takefocus 0 -font {-size 10} 
ttk::label .fiLes.info.2.3.frequencyX -text "0" -takefocus 0 -font {-size 10} -width 10 
ttk::label .fiLes.info.2.chst -text "" -takefocus 0 -font {-size 10} -width 10 
#
ttk::label .fiLes.fiLes1.name -text $lang_var(.fiLes.fiLes1.name) -takefocus 0 -font {-size 10} -width 9
ttk::label .fiLes.fiLes2.name -text $lang_var(.fiLes.fiLes2.name) -takefocus 0 -font {-size 10} -width 9
entry .fiLes.fiLes1.textFile -width 54 -textvariable inpF
.fiLes.fiLes1.textFile insert end $dirIn
.fiLes.fiLes1.textFile xview [string lengt $dirIn]
entry .fiLes.fiLes2.textFile -width 54 -textvariable outF
.fiLes.fiLes2.textFile insert end $dirOut
.fiLes.fiLes2.textFile  xview [string lengt $dirOut]
ttk::button .fiLes.fiLes1.inputFile -command inFile -width 2 -text ...
ttk::button .fiLes.fiLes2.saveFile -command sFile -width 2 -text ...
# Section DVD
frame .fiLes.vframe -bg #b4d9cf
checkbutton .fiLes.vframe.dvdvideo -variable dvd -text $lang_var(.fiLes.vframe.dvdvideo) -width 17 -anchor w -command {dvdcheck $dvd} -bg #b4d9cf
ttk::label .fiLes.vframe.ltemp -text $lang_var(.fiLes.vframe.ltemp) -state disabled -width 16 -anchor e -takefocus 0 -font {-size 10} -background #b4d9cf
entry .fiLes.vframe.temp -width 28 -textvariable tempdir -state disabled
ttk::button .fiLes.vframe.tempfolder -command TempFolder -width 2 -state disabled -text ... 
.fiLes.vframe.temp xview [string lengt $tempdir]
# Section packed
frame .fiLes.packed -bg #b4d9cf
checkbutton .fiLes.packed.dvdvideo -variable multi -text $lang_var(.fiLes.packed.dvdvideo) -width 17 -anchor w -command {multicheck $multi} -bg #b4d9cf
ttk::label .fiLes.packed.ltemp -text $lang_var(.fiLes.packed.ltemp) -state disabled -width 16 -anchor e -takefocus 0 -font {-size 10} -background #b4d9cf
entry .fiLes.packed.temp -width 28 -textvariable multiple -state disabled
ttk::button .fiLes.packed.tempfolder -command MultiFolder -width 2 -state disabled -text ... 
.fiLes.packed.temp xview [string lengt $multiple]
# Section pack
pack .fiLes -fill x -pady 10 -padx 5
pack .fiLes.fiLes1 -fill x
pack .fiLes.fiLes1.name -side left -padx 5 
pack .fiLes.fiLes1.textFile -side left -padx 5
pack .fiLes.fiLes1.inputFile -side left -padx {0 5}
#
pack .fiLes.info -fill x -pady {5 5} -padx {85 5}
pack .fiLes.info.1 -side left -padx 5 -pady 2 -fill both -expand 1
pack .fiLes.info.1.1 -side top -anchor nw
pack .fiLes.info.1.1.vcodec -side left
pack .fiLes.info.1.1.vcodecX -side left
pack .fiLes.info.1.2 -side top -anchor nw
pack .fiLes.info.1.2.vbit -side left
pack .fiLes.info.1.2.vbitX -side left
pack .fiLes.info.1.3 -side top -anchor nw
pack .fiLes.info.1.3.size -side left
pack .fiLes.info.1.3.sizeX -side left
pack .fiLes.info.1.4 -side top -anchor nw
pack .fiLes.info.1.4.aspect -side left
pack .fiLes.info.1.4.aspectX -side left
pack .fiLes.info.2 -side left -padx 5 -pady 2 -fill both -expand 1
pack .fiLes.info.2.1 -side top  -anchor nw
pack .fiLes.info.2.1.acodec -side left
pack .fiLes.info.2.1.acodecX -side left
pack .fiLes.info.2.2 -side top  -anchor nw
pack .fiLes.info.2.2.abit -side left
pack .fiLes.info.2.2.abitX -side left
pack .fiLes.info.2.3 -side top  -anchor nw
pack .fiLes.info.2.3.frequency -side left
pack .fiLes.info.2.3.frequencyX -side left
pack .fiLes.info.2.chst -side top  -anchor nw
#
pack .fiLes.vframe -fill x -pady {0 5}
pack .fiLes.vframe.dvdvideo -side left
pack .fiLes.vframe.ltemp -side left -padx {0 5}
pack .fiLes.vframe.temp -side left -padx {0 5}
pack .fiLes.vframe.tempfolder -side left -padx {0 5} -pady 2
#
pack .fiLes.packed -fill x -pady {0 5}
pack .fiLes.packed.dvdvideo -side left
pack .fiLes.packed.ltemp -side left -padx {0 5}
pack .fiLes.packed.temp -side left -padx {0 5}
pack .fiLes.packed.tempfolder -side left -padx {0 5} -pady 2
#
pack .fiLes.fiLes2 -fill x -pady {0 5}
pack .fiLes.fiLes2.name -side left -padx 5
pack .fiLes.fiLes2.textFile -side left -padx 5
pack .fiLes.fiLes2.saveFile -side left -padx {0 5}
# Section options
ttk::frame .option
ttk::labelframe .option.optionV -text $lang_var(.option.optionV)
ttk::labelframe .option.optionA -text $lang_var(.option.optionA)

pack .option -fill x
pack .option.optionV -side left -padx 5 -pady 2 -fill both -expand 1
pack .option.optionA -side left -padx {0 5} -pady 2 -fill both -expand 1

#
ttk::frame .option.optionV.v1
pack .option.optionV.v1 -side left -pady 5
ttk::frame .option.optionV.v1.1
pack .option.optionV.v1.1 -fill x -padx 5
spinbox .option.optionV.v1.1.sbitrate -from 50 -to 10000 -width 7 -textvariable vBit
ttk::label .option.optionV.v1.1.bitrate -text $lang_var(.option.optionV.v1.1.bitrate) -takefocus 0 -font {-size 10}
pack .option.optionV.v1.1.sbitrate -side left
pack .option.optionV.v1.1.bitrate -side left -padx {5 0}
ttk::frame .option.optionV.v1.2  
pack .option.optionV.v1.2 -pady 5 -fill x -padx 5
spinbox .option.optionV.v1.2.sfps -from 5.0 -to 35.0 -increment 0.5 -width 18 -textvariable fps
ttk::label .option.optionV.v1.2.fps -text $lang_var(.option.optionV.v1.2.fps) -takefocus 0 -font {-size 10}
pack .option.optionV.v1.2.sfps -side left
pack .option.optionV.v1.2.fps -side left -padx {5 0}
ttk::frame .option.optionV.v1.3
pack .option.optionV.v1.3 -fill x -padx 5 -pady {0 5}
ttk::checkbutton .option.optionV.v1.3.chdein -variable dIn -text $lang_var(.option.optionV.v1.3.chdein)
pack .option.optionV.v1.3.chdein -side left
ttk::frame .option.optionV.v1.4
pack .option.optionV.v1.4 -fill x -pady {0 5} -padx 5
ttk::combobox .option.optionV.v1.4.lst -values $vcodecS -width 18 -textvariable vcodec
if {[lsearch -exact $vcodecS $vcodec]!=-1} {
	.option.optionV.v1.4.lst current [lsearch -exact $vcodecS $vcodec]
} else {.option.optionV.v1.4.lst current 1}
ttk::label .option.optionV.v1.4.label -text $lang_var(.option.optionV.v1.4.label) -takefocus 0 -font {-size 10}
pack .option.optionV.v1.4.lst -side left
pack .option.optionV.v1.4.label -side left -padx {5 0}
ttk::frame .option.optionV.v1.5
pack .option.optionV.v1.5 -fill x -pady {0 5} -padx 5
ttk::checkbutton .option.optionV.v1.5.copy -variable vcopy -command {vcopy $vcopy} -text $lang_var(.option.optionV.v1.5.copy)
pack .option.optionV.v1.5.copy -side left
#
ttk::labelframe .option.optionV.size -text $lang_var(.option.optionV.size)
pack .option.optionV.size -side left -padx 5 -anchor nw
ttk::frame .option.optionV.size.2  
pack .option.optionV.size.2 -pady 5 -fill x -padx 10
spinbox .option.optionV.size.2.sx -from 10 -to 1024 -increment 2 -width 5 -textvariable X
ttk::label .option.optionV.size.2.x1 -text "X" -takefocus 0 -font {-size 10}
pack .option.optionV.size.2.sx -side left
pack .option.optionV.size.2.x1 -side left -padx {5 0}
ttk::frame .option.optionV.size.3  
pack .option.optionV.size.3 -pady 5 -fill x -padx 10
spinbox .option.optionV.size.3.sy -from 10 -to 1024 -increment 2 -width 5 -textvariable Y
ttk::label .option.optionV.size.3.y -text "Y" -takefocus 0 -font {-size 10}
pack .option.optionV.size.3.sy -side left
pack .option.optionV.size.3.y -side left -padx {5 0}
#
ttk::frame .option.optionA.1 
pack .option.optionA.1 -padx 10 -pady {5 0} -fill x
ttk::label .option.optionA.1.ac -text $lang_var(.option.optionA.1.ac) -takefocus 0 -font {-size 10}
ttk::combobox .option.optionA.1.lstac -values $acodecS -width 5 -textvariable acodec
if {[lsearch -exact $acodecS $acodec]!=-1} {
	.option.optionA.1.lstac current [lsearch -exact $acodecS $acodec]
} else {.option.optionA.1.lstac current 0}
pack .option.optionA.1.lstac -side left
pack .option.optionA.1.ac -side left -padx {5 0}
ttk::frame .option.optionA.2 
pack .option.optionA.2 -padx 10 -pady {10 0} -fill x
ttk::label .option.optionA.2.fr -text $lang_var(.option.optionA.2.fr) -takefocus 0 -font {-size 10}
ttk::combobox .option.optionA.2.lstfr -values $freqS -width 5 -textvariable freq
if {[lsearch -exact $freqS $freq]!=-1} {
	.option.optionA.2.lstfr current [lsearch -exact $freqS $freq]
} else {.option.optionA.2.lstfr current 1}
pack .option.optionA.2.lstfr -side left
pack .option.optionA.2.fr -side left -padx {5 0}
ttk::frame .option.optionA.3 
pack .option.optionA.3 -padx 10 -pady {10 0} -fill x
ttk::label .option.optionA.3.bit -text $lang_var(.option.optionA.3.bit) -takefocus 0 -font {-size 10}
ttk::combobox .option.optionA.3.lstbit -values $aBitS -width 5 -textvariable aBit
if {[lsearch -exact $aBitS $aBit]!=-1} {
	.option.optionA.3.lstbit current [lsearch -exact $aBitS $aBit]
} else {.option.optionA.3.lstbit current 3}
pack .option.optionA.3.lstbit -side left
pack .option.optionA.3.bit -side left -padx {5 0}
ttk::frame .option.optionA.4 
pack .option.optionA.4 -padx 10 -pady {5 3} -fill x
ttk::checkbutton .option.optionA.4.chst -variable stu -text $lang_var(.option.optionA.4.chst)
pack .option.optionA.4.chst -side left
ttk::frame .option.optionA.5
pack .option.optionA.5 -fill x -pady {0 5} -padx 10
ttk::checkbutton .option.optionA.5.copy -variable acopy -command {acopy $acopy} -text $lang_var(.option.optionA.5.copy)
pack .option.optionA.5.copy -side left
#
ttk::frame .addit
pack .addit -fill x
ttk::label .addit.moreL -text $lang_var(.addit.moreL) 
entry .addit.more -textvariable more_opt -width 120
pack .addit.moreL -side left -padx 5 -pady 2
pack .addit.more  -side left -padx {0 5} -pady 2
#
frame .but -bg #b4d9cf
checkbutton .but.chop -variable openFu -text $lang_var(.but.chop) -bg #b4d9cf
ttk::button .but.go -command {.log.msg insert end [go $inpF $outF $vBit $fps $dIn $X $Y $vcodec $acodec $freq $aBit $stu $more_opt] red_color} -text $lang_var(.but.go) -width 7
ttk::button .but.cancel -command {save_options} -text $lang_var(.but.cancel) -width 7
ttk::button .but.option -command {adv_opt_window} -text $lang_var(.but.option) -width 10
pack .but -fill x
pack .but.chop -side left -padx {5 5}
pack .but.cancel -side right -padx {5 5} -pady 2
pack .but.go -side right -padx {5 5}
pack .but.option -side right -padx 5

pack [frame .log] -fill x
text .log.msg -heigh 3 -yscrollcommand {.log.scroll set} 
.log.msg tag configure red_color -foreground red
.log.msg tag configure blue_color -foreground blue -font "bold"
scrollbar .log.scroll  -command {.log.msg yview}
pack .log.msg -fill both -side left
pack .log.scroll -fill y -side right

#END
