PGDMP         1                 y            mydb #   12.5 (Ubuntu 12.5-0ubuntu0.20.10.1) #   12.5 (Ubuntu 12.5-0ubuntu0.20.10.1)     N           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            O           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            P           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            Q           1262    16385    mydb    DATABASE     v   CREATE DATABASE mydb WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';
    DROP DATABASE mydb;
                postgres    false            �            1259    16394    option    TABLE     �   CREATE TABLE public.option (
    option_text character varying(100)[] NOT NULL,
    question_id integer NOT NULL,
    votes integer DEFAULT 0 NOT NULL,
    option_id integer NOT NULL
);
    DROP TABLE public.option;
       public         heap    postgres    false            �            1259    16467    option_option_id_seq    SEQUENCE     �   CREATE SEQUENCE public.option_option_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 +   DROP SEQUENCE public.option_option_id_seq;
       public          postgres    false    202            R           0    0    option_option_id_seq    SEQUENCE OWNED BY     M   ALTER SEQUENCE public.option_option_id_seq OWNED BY public.option.option_id;
          public          postgres    false    207            �            1259    16402    question    TABLE     �   CREATE TABLE public.question (
    question_id integer NOT NULL,
    question_text character varying(200)[] NOT NULL,
    question_system integer DEFAULT 0 NOT NULL
);
    DROP TABLE public.question;
       public         heap    postgres    false            �            1259    16408    question_question_id_seq    SEQUENCE     �   CREATE SEQUENCE public.question_question_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 /   DROP SEQUENCE public.question_question_id_seq;
       public          postgres    false    203            S           0    0    question_question_id_seq    SEQUENCE OWNED BY     U   ALTER SEQUENCE public.question_question_id_seq OWNED BY public.question.question_id;
          public          postgres    false    204            �            1259    16410    userkey    TABLE     �   CREATE TABLE public.userkey (
    user_key character(32)[] NOT NULL,
    question_id integer NOT NULL,
    user_id integer NOT NULL,
    used boolean DEFAULT false NOT NULL
);
    DROP TABLE public.userkey;
       public         heap    postgres    false            �            1259    16416    user_user_id_seq    SEQUENCE     �   CREATE SEQUENCE public.user_user_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 '   DROP SEQUENCE public.user_user_id_seq;
       public          postgres    false    205            T           0    0    user_user_id_seq    SEQUENCE OWNED BY     H   ALTER SEQUENCE public.user_user_id_seq OWNED BY public.userkey.user_id;
          public          postgres    false    206            �           2604    16469    option option_id    DEFAULT     t   ALTER TABLE ONLY public.option ALTER COLUMN option_id SET DEFAULT nextval('public.option_option_id_seq'::regclass);
 ?   ALTER TABLE public.option ALTER COLUMN option_id DROP DEFAULT;
       public          postgres    false    207    202            �           2604    16420    question question_id    DEFAULT     |   ALTER TABLE ONLY public.question ALTER COLUMN question_id SET DEFAULT nextval('public.question_question_id_seq'::regclass);
 C   ALTER TABLE public.question ALTER COLUMN question_id DROP DEFAULT;
       public          postgres    false    204    203            �           2604    16421    userkey user_id    DEFAULT     o   ALTER TABLE ONLY public.userkey ALTER COLUMN user_id SET DEFAULT nextval('public.user_user_id_seq'::regclass);
 >   ALTER TABLE public.userkey ALTER COLUMN user_id DROP DEFAULT;
       public          postgres    false    206    205            F          0    16394    option 
   TABLE DATA           L   COPY public.option (option_text, question_id, votes, option_id) FROM stdin;
    public          postgres    false    202   �       G          0    16402    question 
   TABLE DATA           O   COPY public.question (question_id, question_text, question_system) FROM stdin;
    public          postgres    false    203   4       I          0    16410    userkey 
   TABLE DATA           G   COPY public.userkey (user_key, question_id, user_id, used) FROM stdin;
    public          postgres    false    205   �       U           0    0    option_option_id_seq    SEQUENCE SET     B   SELECT pg_catalog.setval('public.option_option_id_seq', 2, true);
          public          postgres    false    207            V           0    0    question_question_id_seq    SEQUENCE SET     F   SELECT pg_catalog.setval('public.question_question_id_seq', 2, true);
          public          postgres    false    204            W           0    0    user_user_id_seq    SEQUENCE SET     >   SELECT pg_catalog.setval('public.user_user_id_seq', 4, true);
          public          postgres    false    206            �           2606    16471    option option_pkey 
   CONSTRAINT     W   ALTER TABLE ONLY public.option
    ADD CONSTRAINT option_pkey PRIMARY KEY (option_id);
 <   ALTER TABLE ONLY public.option DROP CONSTRAINT option_pkey;
       public            postgres    false    202            �           2606    16427    question question_pkey 
   CONSTRAINT     ]   ALTER TABLE ONLY public.question
    ADD CONSTRAINT question_pkey PRIMARY KEY (question_id);
 @   ALTER TABLE ONLY public.question DROP CONSTRAINT question_pkey;
       public            postgres    false    203            �           2606    16429    userkey user_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.userkey
    ADD CONSTRAINT user_pkey PRIMARY KEY (user_id);
 ;   ALTER TABLE ONLY public.userkey DROP CONSTRAINT user_pkey;
       public            postgres    false    205            �           2606    16435    option question_id    FK CONSTRAINT     �   ALTER TABLE ONLY public.option
    ADD CONSTRAINT question_id FOREIGN KEY (question_id) REFERENCES public.question(question_id) NOT VALID;
 <   ALTER TABLE ONLY public.option DROP CONSTRAINT question_id;
       public          postgres    false    3011    202    203            �           2606    16440    userkey question_id    FK CONSTRAINT     �   ALTER TABLE ONLY public.userkey
    ADD CONSTRAINT question_id FOREIGN KEY (question_id) REFERENCES public.question(question_id) NOT VALID;
 =   ALTER TABLE ONLY public.userkey DROP CONSTRAINT question_id;
       public          postgres    false    203    3011    205            F   �   x�-�A�0����cW�c�.����4��R+��\�w�x/Y�/y+�:Xv���@�V��bWra�����oǺ�A
k،�!o�ik�޸����6hri�T�>v��L+����NT-&�����|�X�I�%O�r6
�p�OM��qX	!~^�<_      G   g   x�3�Vr��WH�IM�W��L�W��4�2
�'�&��*T�%&g'*$��$%�����&@I����{�r���K�R��rR�.<�bA� ɨ(�      I   c   x��6
�H5/s6�I+3)5�L��.5Έ�57u51��4�4�L�v$�������`�&�f ��D(4)4$B�)Han4)t&B�1Pa� ��Uc     