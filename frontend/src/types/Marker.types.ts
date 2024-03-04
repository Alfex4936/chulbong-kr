export interface Photo {
  photoId: number;
  markerId: number;
  photoUrl: string;
  uploadedAt: Date;
}

export interface Marker {
  markerId: number;
  userId: number;
  latitude: number;
  longitude: number;
  description: string;
  createdAt: Date;
  updatedAt: Date;
  username: string;
  photos?: Photo[];
  dislikeCount?: number;
}
