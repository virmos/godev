import useSWR from "swr"
import axios from "axios";

export const handler = () => (productId) => {
    const swrRes = useSWR(`/api/admin/product`,
        async () => {
            const { data } = await axios.get(`/api/admin/products/${productId}`);
            return data
        }
    )
  
    return swrRes
}