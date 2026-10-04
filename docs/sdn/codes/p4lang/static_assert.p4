#include <core.p4>
#include <v1model.p4>

const bool _check = static_assert(V1MODEL_VERSION > 20180000,
                                  "Expected V1MODEL_VERSION > 20180000");
