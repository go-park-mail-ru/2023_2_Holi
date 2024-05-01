from apscheduler.schedulers.blocking import BlockingScheduler
from get_data import get_ratings
from logic import count_recommendations
from save_data import save_recommendations

def job():
    data = get_ratings()
    top_10 = count_recommendations(data)
    save_recommendations(top_10)

scheduler = BlockingScheduler()
scheduler.add_job(job, 'interval', minutes=1)
scheduler.start()