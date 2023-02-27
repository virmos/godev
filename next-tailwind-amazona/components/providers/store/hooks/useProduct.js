import useSWR from "swr"
import axios from "axios";

export const handler = () => (productId) => {
    const swrRes = useSWR(`/api/products/${productId}`,
        async () => {
            const { data } = await axios.get(`/api/products/${productId}`);
            return data
        }
    )
  
    return swrRes
}