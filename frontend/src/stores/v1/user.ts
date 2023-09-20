import { create } from "zustand";
import * as api from "../../helpers/api";
import { UserData, User, UserId } from "../../types/user";


const convertResponseModelUser = (data: UserData): User => {
  return {
    ...data.user,
  };
};

interface UserState {
	userMapById: Record<UserId, User>;
	currentUserId?: UserId;
	fetchCurrentUser: () => Promise<User>;
	getCurrentUser: () => User;
}

const useUserStore = create<UserState>()((set, get) => ({
	userMapById: {},
	fetchCurrentUser: async () => {
		const { data } = await api.getCurrentUser();
		const user = convertResponseModelUser(data);
		const userMap = get().userMapById;
		userMap[user.id] = user;
		set({userMapById: userMap, currentUserId: user.id});
		return user;
	},
	getCurrentUser: () => {
		const userMap = get().userMapById;
		const currentUserId = get().currentUserId;
		return userMap[currentUserId as string];
	},
}));

export default useUserStore;