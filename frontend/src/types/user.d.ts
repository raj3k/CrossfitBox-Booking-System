interface UserProfile {
  id: string;
  user_id: string;
  phone_number: string | null;
  birth_date: string | null;
}

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  is_active: boolean;
  is_staff: boolean;
  is_superuser: boolean;
  thumbnail: string | null;
  created_at: string;
  profile: UserProfile;
}

export interface UserData {
  user: User;
}