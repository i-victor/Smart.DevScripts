
// valac --pkg librsvg-2.0 --pkg gtk+-3.0 sample2.vala

class PiecePreviewWindow : Gtk.Window
{
    public PiecePreviewWindow( ) {
        set_title( "Piece Preview" );
        add( new PieceView( ) );
        show_all( );
    }

    public override void destroy( ) {
        Gtk.main_quit( );
    }

    public static int main( string[] args ) {
        Gtk.init( ref args );
        new PiecePreviewWindow( );
        Gtk.main( );

        return 0;
    }
}

class PieceView : Gtk.DrawingArea
{
    public PieceView( ) {
        set_size_request( 400, 400 );
        add_events( Gdk.EventMask.BUTTON_PRESS_MASK | Gdk.EventMask.SCROLL_MASK );

        temp_surface = new Cairo.ImageSurface( Cairo.Format.ARGB32, 12 * square_size, square_size );

        string format = "./sample%s.svg";
        render_piece( format.printf("2"), 0 );

    }

    public override bool draw( Cairo.Context cr ) {
        cr.set_source_rgb( 0.5, 0.5, 0.5 );
        cr.rectangle( 0.0, 0.0, 400.0, 400.0 );
        cr.fill( );

        cr.scale( 400.0 / square_size, 400.0 / square_size );
        cr.set_source_surface( temp_surface, -offset * square_size, 0 );
        cr.rectangle( 0, 0, square_size, square_size );
        cr.clip( );
        cr.paint( );

        return false;
    }

    void render_piece( string file_name, int offset ) {
        Rsvg.Handle handle;
        try {
            handle = new Rsvg.Handle.from_file( file_name );
        } catch( Error e ) {
            stderr.printf( "can not open svg file\n" );
            return;
        }

        temp_cr = new Cairo.Context( temp_surface );
        temp_cr.save( );
        temp_cr.translate( square_size * offset, 0 );
        temp_cr.scale( (double) square_size / handle.width, (double) square_size / handle.height );
        handle.render_cairo( temp_cr );
        temp_cr.restore( );
    }

    public override bool button_press_event( Gdk.EventButton event ) {
        if( event.button  == 1 ) {
            change_piece( true );
        } else if( event.button == 3 ) {
            change_piece( false );
        } else
            return false;

        return true;
    }

    public override bool scroll_event( Gdk.EventScroll event ) {
        if( event.direction == Gdk.ScrollDirection.UP ) {
            change_piece( true );
        } else if( event.direction == Gdk.ScrollDirection.DOWN ) {
            change_piece( false );
        } else
            return false;
        return true;
    }

    private void change_piece( bool forward ) {
        if( forward ) offset ++;
        else offset --;

        if( offset < 0 ) offset = 11;
        if( offset > 11) offset = 0;

        queue_draw( );
    }

    private Cairo.ImageSurface temp_surface;
    private Cairo.Context temp_cr;
    private int square_size = 500;
    private int offset;
}

// #END
