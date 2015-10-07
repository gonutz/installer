//----------------------------------------------------------------------------
// registry.h
//
// registry manipulation stuff
//
// Copyright (C) 2011 Neil Butterworth
//----------------------------------------------------------------------------

#ifndef INV_PATHED_REGISTRY_H
#define INV_PATHED_REGISTRY_H

#include <windows.h>
#include <string>
#include <vector>

class RegPath {

	public:

		enum Env { Expand, NoExpand };

		RegPath( HKEY root );
		~RegPath();

		unsigned int Count() const;
		std::string At( unsigned int  i ) const;
		std::string Str() const;

		bool Find( const std::string & apath,  Env ev ) const;
		bool Add( const std::string & apath );
		bool Remove( const std::string & apath );
		void ReplaceAll( const std::string & apath );

	private:

		void SplitPath( const std::string & path );
		void UpdateReg();
		bool  RemoveIC( const std::string & adir );

		HKEY mRoot, mPathKey;
		typedef std::vector <std::string>  VecType;
		VecType mPath;
};


#endif
