import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/router';
import React from 'react';
import { toast } from 'react-toastify';
import Product from '@models/Product';
import db from '@utils/db';
import { useStore } from '@components/providers';
import { BaseLayout } from '@components/ui/layout';
import { getProduct } from '@components/api';
import { Curriculum } from '@components/ui/course';

export default function ProductScreen(props) {
    const { product } = props;
    const { state, dispatch } = useStore();
    const router = useRouter();
    if (!product) {
        return <>Product Not Found</>;
    }

    const addToCartHandler = async () => {
        const existItem = state.cart.cartItems.find((x) => x.slug === product.slug);
        const quantity = existItem ? existItem.quantity + 1 : 1;
        const { data } = await getProduct(product._id);

        if (data.countInStock < quantity) {
            return toast.error('Sorry. Product is out of stock');
        }

        dispatch({ type: 'CART_ADD_ITEM', payload: { ...product, quantity } });
        router.push('/cart');
    };

    return (
        <>
            <div className="relative bg-white overflow-hidden">
                <div className="flex justify-end">
                    <div className="py-2">
                        <Link href="/">Back to products</Link>
                    </div>
                </div>
                <div className="max-w-7xl mx-auto flex flex-col sm:flex-row">
                    <div className="flex-1 relative z-10 pb-8 bg-white sm:pb-16 md:pb-20 lg:max-w-2xl lg:w-full lg:pb-28 xl:pb-32">
                        <svg className="hidden lg:block absolute right-0 inset-y-0 h-full w-48 text-white transform translate-x-1/2" fill="currentColor" viewBox="0 0 100 100" preserveAspectRatio="none" aria-hidden="true">
                            <polygon points="50,0 100,0 50,100 0,100" />
                        </svg>

                        <main className="mt-2 mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                            <div className="text-center lg:text-left">
                                <div className="text-xl inline-block p-4 py-2 rounded-full font-bold bg-green-200 text-green-700">
                                    Brand new:
                                </div>
                                <h1 className="text-4xl tracking-tight font-extrabold text-gray-900 sm:text-5xl md:text-6xl">
                                    <span className="block lg:inline">
                                        {product.name.substring(0, product.name.length / 2)}
                                    </span>
                                    <span className="block text-indigo-600 lg:inline">
                                        {product.name.substring(product.name.length / 2)}
                                    </span>
                                </h1>
                                <p className="mt-3 text-base text-gray-500 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-0">
                                    {product.description}
                                </p>
                                <div className="card p-5 mt-5 w-52 mx-auto">
                                    <div className="mb-2 flex justify-between">
                                        <div>Price</div>
                                        <div>${product.price}</div>
                                    </div>
                                    <div className="mb-2 flex justify-between">
                                        <div>Status</div>
                                        <div>{product.countInStock > 0 ? 'In stock' : 'Unavailable'}</div>
                                    </div>
                                    <button
                                        className="primary-button w-full"
                                        onClick={addToCartHandler}
                                    >
                                        Add to cart
                                    </button>
                                </div>
                            </div>
                        </main>
                    </div>
                    {/* <div className="flex-1 lg:absolute lg:inset-y-0 lg:right-0 lg:w-2/3"> */}
                    <div className="flex-1">
                        <Image
                            className="h-56 w-full object-cover sm:h-72 md:h-96 lg:w-full lg:h-full"
                            src={product.image}
                            alt={product.name}
                            width={640}
                            height={640}
                            layout="responsive"
                        />
                    </div>
                </div>
            </div>
            <Curriculum points={[
                {key: 'Name', value: product.name},
                {key: 'Category', value: product.category},
                {key: 'Brand', value: product.brand},
                {key: 'Rating', value: `${product.rating} of ${product.numReviews} reviews`},
            ]}/>
        </>
    );
}

export async function getServerSideProps(context) {
    const { params } = context;
    const { slug } = params;

    await db.connect();
    const product = await Product.findOne({ slug }).lean();
    await db.disconnect();
    return {
        props: {
            product: product ? db.convertDocToObj(product) : null,
        },
    };
}

ProductScreen.Layout = BaseLayout