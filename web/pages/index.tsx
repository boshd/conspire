import type { NextPage } from 'next'
import Head from 'next/head'
import Image from 'next/image'
import styles from '../styles/Home.module.css'
import { Form } from './Form'

const Home: NextPage = () => {
  return (
    <div className={styles.container}>
      <Head>
        <title>Conspiracy theories app u/ cohere</title>
      </Head>

      <main className={styles.main}>
        <Form></Form>
      </main>
    </div>
  )
}

export default Home
