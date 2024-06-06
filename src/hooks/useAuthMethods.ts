import { useQuery } from "@tanstack/react-query"
import { usePocket } from "../contexts/pocketbase"

export const useAuthMethods = () => {
    const { pb } = usePocket();

    return useQuery({
        queryKey: ['authMethods'],
        queryFn: () => {
            return pb?.collection("users").listAuthMethods()
        },
    })
}