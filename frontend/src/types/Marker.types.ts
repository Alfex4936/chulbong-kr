export interface Photo {
  photoId: number;
  markerId: number;
  photoUrl: string;
  uploadedAt: Date;
}

export interface Marker {
  markerId: number;
  userId: number | null;
  latitude: number;
  longitude: number;
  description: string;
  createdAt: Date;
  updatedAt: Date;
  username: string;
  photos?: Photo[];
  dislikeCount?: number;
  disliked: boolean;
  addr?: string;
  isChulbong?: boolean;
  favorited?: boolean;
  address?: string;
  favCount?:number;
}
