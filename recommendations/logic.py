from surprise import Dataset, Reader
from collections import defaultdict
from surprise import SVD


def count_recommendations(ratings):
    reader = Reader(rating_scale=(0.5, 10))
    data = Dataset.load_from_df(ratings[['user_id', 'video_id', 'rate']], reader)

    # Обучаем модель на всех доступных данных
    trainset = data.build_full_trainset()
    algo = SVD(n_factors=100, n_epochs=10, lr_all=0.005, reg_all=0.02)
    algo.fit(trainset)

    # Теперь нам нужно пройтись по всем пользователям и фильмам и сделать предсказания для ненаблюдаемых пар
    testset = trainset.build_anti_testset()

    predictions = algo.test(testset)
    # Получаем топ-N рекомендаций для каждого пользователя
    top_n = get_top_n(predictions, n=10)
    # Выводим рекомендации для каждого пользователя
    for uid, user_ratings in top_n.items():
        print(uid, [iid for (iid, _) in user_ratings])

    return top_n


def get_top_n(predictions, n=10):
    # Словарь для хранения рекомендаций для каждого пользователя
    top_n = defaultdict(list)

    for uid, iid, true_r, est, _ in predictions:
        top_n[uid].append((iid, est))

    # Сортировка прогнозов для каждого пользователя и выборка N наибольших
    for uid, user_ratings in top_n.items():
        user_ratings.sort(key=lambda x: x[1], reverse=True)
        top_n[uid] = user_ratings[:n]

    return top_n
