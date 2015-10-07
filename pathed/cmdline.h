//----------------------------------------------------------------------------
// cmdline.h
//
// Simple command line handler class, supporting shifting of parameters. The
// name of the command i.e. argv[0] is preserved.
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------

#ifndef INC_PATHED_CMDLINE_H
#define INC_PATHED_CMDLINE_H

#include <vector>
#include <string>

class CmdLine {

	public:

		CmdLine( int argc, char * argv[] ) {
			for ( int i = 0; i < argc; i++ ) {
				mArgv.push_back( argv[i] );
			}
		}

		int Argc() const {
			return mArgv.size();
		}

		std::string Argv( int i ) const {
			return mArgv.at( i );
		}

		void Shift( int n = 1 ) {
			while( n-- && Argc() > 1 ) {
				mArgv.erase( mArgv.begin() + 1 );
			}
		}

	private:

		std::vector <std::string> mArgv;
};

#endif

