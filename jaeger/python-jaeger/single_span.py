import requests
import logging
import time
from random import randint
from jaeger_client import Config


def download():
    rsp = requests.get("https://www.liwenzhou.com")
    return rsp


def parser():
    time.sleep(randint(1, 9) * 0.1)

# 如果这个函数有多个操作，我们希望看到每个函数的执行时间


def insert_to_mysql(parent_span):
    with tracer.start_span("select", child_of=parent_span) as select_span:
        # 数据库的执行时间
        time.sleep(randint(1, 9) * 0.1)

    with tracer.start_span("exec", child_of=parent_span) as exec_span:
        # 插入数据的执行时间
        time.sleep(randint(1, 9) * 0.1)

if __name__ == "__main__":
    # 前边的这部分属于是日志
    log_level = logging.DEBUG
    logging.getLogger('').handlers = []
    logging.basicConfig(format='%(asctime)s %(message)s', level=log_level)
    # config这里是配置
    config = Config(
        config={  # usually read from some yaml config
            'sampler': {
                'type': 'const',  # 全部采用
                'param': 1,  # 1.代表开启全部采样
            },
            'local_agent': {
                'reporting_host': '127.21.17.5',
                'reporting_port': '6831',
            },
            'logging': True,
        },
        service_name='Hello',  # 链路的名称
        validate=True,
    )
    # this call also sets opentracing.tracer
    tracer = config.initialize_tracer()

    with tracer.start_span("spider") as spider_span:

        with tracer.start_span("get", child_of=spider_span) as get_span:
            download()
    # 一个span监控一个服务，span服务内的视图
    # 使用同一个trace生成的span不放在一起
    # 怎么解决 让多个span作为一个父span的子span

        with tracer.start_span("parser", child_of=spider_span) as parser_span:
            parser()

        with tracer.start_span("insert", child_of=spider_span) as insert_span:
            insert_to_mysql(insert_span)
    # yield to IOLoop to flush the spans -
    # https://github.com/jaegertracing/jaeger-client-python/issues/50
    time.sleep(2)
    tracer.close()  # flush any buffered spans
