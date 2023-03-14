import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/admin/products`,
        async () => {
            const { data } = await axios.get(`/api/admin/products`);
            return data
        }
    )
  
    return swrRes
}