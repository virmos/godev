import axios from 'axios';
import { toast } from 'react-toastify';
import { BaseLayout } from '@components/ui/layout';
import ProductItem from '@components/ui/product/item/ProductItem';
import Product from '@models/Product';
import db from '@utils/db';
import { useStore } from '@components/providers/store';
import { Carousel } from 'react-responsive-carousel';
import 'react-responsive-carousel/lib/styles/carousel.min.css';
import Link from 'next/link';
        
export default function Home({ products, featuredProducts }) {
    const { state, dispatch } = useStore();
    const { cart } = state;

    const addToCartHandler = async (product) => {
        const existItem = cart.cartItems.find((x) => x.slug === product.slug);
        const quantity = existItem ? existItem.quantity + 1 : 1;
        const { data } = await axios.get(`/api/products/${product._id}`);

        if (data.countInStock < quantity) {
            return toast.error('Sorry. Product is out of stock');
        }
        dispatch({ type: 'CART_ADD_ITEM', payload: { ...product, quantity } });

        toast.success('Product added to the cart');
    };

    return (
        <>
            <Carousel showThumbs={false} autoPlay>
                {featuredProducts.map((product) => (
                    <div key={product._id}>
                        <Link href={`/product/${product.slug}`} passHref>
                            <a className="flex">
                                <img src={product.banner} alt={product.name} />
                            </a>
                        </Link>
                    </div>
                ))}
            </Carousel>
            <h2 className="h2 my-4">Latest Products</h2>
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3 lg:grid-cols-4">
                {products.map((product) => (
                    <ProductItem
                        product={product}
                        key={product.slug}
                        addToCartHandler={addToCartHandler}
                    ></ProductItem>
                ))}
            </div>
        </>
    );
}

export async function getServerSideProps() {
    await db.connect();
    const products = await Product.find().lean();
    const featuredProducts = await Product.find({ isFeatured: true }).lean();
    return {
        props: {
            featuredProducts: featuredProducts.map(db.convertDocToObj),
            products: products.map(db.convertDocToObj),
        },
    };
}

Home.Layout = BaseLayout