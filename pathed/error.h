//----------------------------------------------------------------------------
// error.h
//
// error reporting for pathed
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------

#ifndef INC_PATHED_ERROR_H
#define INC_PATHED_ERROR_H

#include <windows.h>
#include <exception>
#include <string>

//----------------------------------------------------------------------------
// We throw these kind of exceptions
//----------------------------------------------------------------------------

class Error : public std::exception {

	public:

		Error( const std::string & msg ) : mMsg( msg ) {}
		~Error() throw() {}

		const char *what() const throw() {
			return mMsg.c_str();
		}

	private:
		std::string mMsg;
};

//----------------------------------------------------------------------------
// Get last Windows error as string - code copied from MSDN examples.
//----------------------------------------------------------------------------

inline std::string LastWinError() {
	char * lpMsgBuf;
	::FormatMessage( FORMAT_MESSAGE_ALLOCATE_BUFFER |
					FORMAT_MESSAGE_FROM_SYSTEM |
					FORMAT_MESSAGE_IGNORE_INSERTS,
					NULL,
					::GetLastError(),
					MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
					(LPTSTR) &lpMsgBuf,
					0,
					NULL );
	std::string msg ( lpMsgBuf );
	::LocalFree( lpMsgBuf );
	return  msg;
}
#endif

