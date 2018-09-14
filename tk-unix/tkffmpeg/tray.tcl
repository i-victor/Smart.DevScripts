# tray.tcl

proc balloon {} {
    global progres f inpF burn lang_var
    set wy [winfo screenheight .];
    set wx [winfo screenwidth .];
    set inf [.ic bbox]
    set w .balloon
    catch [destroy $w];
    toplevel $w;
    if {$progres} {
        if {$burn==1} {
            set text "$lang_var(tray_burn)"
        } elseif {$f ne ""} {
            set text "$lang_var(tray_convert) [file tail $f]"
        } elseif {$inpF ne ""} {
            set text "$lang_var(tray_convert) [file tail $inpF]"
        } else {
            set text "tkffmpeg"
        }
        pack [frame $w.prog -background white -relief solid] -fill x
        label $w.prog.lab -text "$text" -background white
        pack $w.prog.lab -fill both -expand true;
        ttk::progressbar $w.prog.p2 -mode indeterminate;
	pack $w.prog.p2 -padx 2 -fill x -expand 1 -pady 2
	$w.prog.p2 start
        
    } else {
        label $w.lab -text "tkffmpeg" -background white -relief solid;
        pack $w.lab -fill both -expand true;
    }
    wm geometry $w +$wx+$wy;
    wm overrideredirect $w true;
    update;
    if {[lindex $inf end]>[expr $wy-[lindex $inf end]]} {
        set y [expr [lindex $inf end]-5]
    } else {
        set y [expr [lindex $inf end]+5]
    }
    if {[lindex $inf 2]>[expr $wx-[winfo width .balloon]-10]} {
        set x [expr [lindex $inf 0]-[winfo width .balloon]-5]
    } else {
        set x [expr [lindex $inf 2]+5]
    }
    wm geometry $w +$x+$y;
    raise $w .;
    focus -force $w;
    after 3000 [list destroy $w];
}
proc ballshow {} {
    set status [wm state .]
    if {$status eq "iconic" || $status eq "withdrawn"} {
        wm deiconify .
    } elseif {$status eq "normal"} {
        wm iconify .
    }
}

#image create photo ico_img -file "$path_proc/tkffmpeg16.png"

tktray::icon .ic -image ico_img;
bind .ic <Enter> balloon;
bind .ic <ButtonPress-1> ballshow;

#END
