type UserId = string;
type ProfileId = string;

interface UserProfile {
  id: ProfileId;
  user_id: UserId;
  phone_number: string | null;
  birth_date: string | null;
}

interface User {
  id: UserId;
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