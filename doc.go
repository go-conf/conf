/*
conf is a way to do configuration for your go app.

What conf does:

    * decodes into your struct
    * can decode from: JSON, commandline flags, env vars
    * Uses pluggable decoders
    * Allows you to customize decoders and priority
    * Makes it easy to write your own decoders for other formats


*/
package conf
