export type Artist = {
  id?: string;
  name: string;
  biography: string;
  genres?: string[];
};

export type Album = {
  id?: string;
  name: string;
  date?: string;
  genre?: string;
  artists?: string[];
};

export type Song = {
  id?: string;
  name: string;
  duration: number;
  album?: string;
  genre?: string;
  artists?: string[];
  audio_url?: string;
};

export type SearchResult = {
  artists: Artist[];
  albums: Album[];
  songs: Song[];
};

export type Genre = {
  id?: string;
  name: string;
  description?: string;
};

export type UserSubscriptions = {
  artists: Artist[];
  genres: Genre[];
};
