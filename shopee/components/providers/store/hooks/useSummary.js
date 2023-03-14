import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/admin/summary`,
        async () => {
            const { data } = await axios.get(`/api/admin/summary`);
            return data
        }
    )
  
    return swrRes
}