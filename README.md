# PWMON

[NAME](#NAME)  
[SYNOPSIS](#SYNOPSIS)  
[DESCRIPTION](#DESCRIPTION)  
[SEE ALSO](#SEE%20ALSO)  
[AUTHOR](#AUTHOR)  

------------------------------------------------------------------------

## NAME <span id="NAME"></span>

pwmon − monitor PipeWire volume and mute status

## SYNOPSIS <span id="SYNOPSIS"></span>

**pwmon**

## DESCRIPTION <span id="DESCRIPTION"></span>

**pwmon** monitors the volume percentage and mute status of the default
PipeWire audio sink. To find the id of the default PipeWire audio sink,
the program parses the output of

**\$ wpctl get-volume @DEFAULT_AUDIO_SINK@**

Using this id, it will then parse the output of

**\$ pw-dump -mN**

and then finds the corresponding sink and reports the volume percentage
and mute status.

## SEE ALSO <span id="SEE ALSO"></span>

*pipewire*(1)

## AUTHOR <span id="AUTHOR"></span>

Kris Andrie Ortega (andrieee44@gmail.com)

------------------------------------------------------------------------
